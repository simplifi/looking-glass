package agent

import (
	"fmt"
	"github.com/jfrog/jfrog-client-go/artifactory"
	"github.com/jfrog/jfrog-client-go/artifactory/auth"
	"github.com/jfrog/jfrog-client-go/artifactory/services"
	aflog "github.com/jfrog/jfrog-client-go/utils/log"
	"github.com/simplifi/looking-glass/config"
	"github.com/simplifi/looking-glass/downloader"
	"log"
	"os"
	"path"
	"time"
)

// Agent monitors a source for changes and pushes files to Artifactory
type Agent struct {
	artifactoryManager artifactory.ArtifactoryServicesManager
	agentDownloader    downloader.Downloader
	agentConfig        config.AgentConfig
	localStoragePath   string
}

// New Agent, pass in the ArtifactoryConfig, and AgentConfig
func New(artifactoryConfig config.ArtifactoryConfig, agentConfig config.AgentConfig) (*Agent, error) {
	artMgr, err := createArtifactoryManager(artifactoryConfig.URL, artifactoryConfig.Key, artifactoryConfig.UserName)
	if err != nil {
		return nil, err
	}
	dl, err := downloader.New(agentConfig.Downloader)
	if err != nil {
		return nil, err
	}
	localStoragePath := path.Join("/tmp", agentConfig.Name)

	agent := Agent{
		artifactoryManager: *artMgr,
		agentDownloader:    dl,
		agentConfig:        agentConfig,
		localStoragePath:   localStoragePath,
	}

	return &agent, nil
}

func createArtifactoryManager(url string, apiKey string, userName string) (*artifactory.ArtifactoryServicesManager, error) {
	// You have to setup a logger for Artifactory client to work
	aflog.SetLogger(aflog.NewLogger(aflog.ERROR, nil))

	details := auth.NewArtifactoryDetails()
	details.SetUrl(url)
	details.SetApiKey(apiKey)
	details.SetUser(userName)

	serviceConfig, err := artifactory.NewConfigBuilder().
		SetArtDetails(details).
		SetDryRun(false).
		Build()
	if err != nil {
		return nil, err
	}

	mgr, err := artifactory.New(&details, serviceConfig)
	if err != nil {
		return nil, err
	}

	return mgr, nil
}

// Start the Agent
func (agt *Agent) Start() {
	for {
		objs, err := agt.agentDownloader.ListObjects()
		if err != nil {
			log.Printf("ERROR: Failed to list objects - %s", err)
		}
		for _, obj := range objs {

			// see if any objects are missing from Artifactory
			if !agt.existsInArtifactory(obj) {
				log.Printf("INFO: [mirror] %s", obj)

				// download object to local storage
				localFile := path.Join(agt.localStoragePath, obj)
				dlErr := agt.agentDownloader.GetObject(obj, localFile)
				if dlErr != nil {
					log.Printf("ERROR: Failed to download object - %v", dlErr)
				}

				// upload to artifactory
				rtErr := agt.uploadToArtifactory(localFile, obj)
				if rtErr != nil {
					log.Printf("ERROR: Failed to upload to Artifactory - %v", rtErr)
				}

				// clean up temp storage
				rmErr := os.RemoveAll(agt.localStoragePath)
				if rmErr != nil {
					log.Printf("ERROR: Failed to clean up temp storage - %v", rmErr)
				}
			} else {
				log.Printf("INFO: [skip] %s", obj)
			}
		}

		log.Printf("INFO: Sleeping for %d seconds", agt.agentConfig.SleepDuration)
		time.Sleep(time.Duration(agt.agentConfig.SleepDuration) * time.Second)
	}
}

func (agt *Agent) uploadToArtifactory(sourceFile string, targetPath string) error {
	params := services.NewUploadParams()
	params.Pattern = sourceFile
	params.Target = fmt.Sprintf("%s/%s", agt.agentConfig.ArtifactoryRepo, targetPath)

	_, _, totalFailed, err := agt.artifactoryManager.UploadFiles(params)

	if err != nil || totalFailed > 0 {
		return fmt.Errorf("ERROR: failed to upload file %q, %v", sourceFile, err)
	}

	return nil
}

func (agt *Agent) existsInArtifactory(filename string) bool {
	exists := false
	params := services.NewSearchParams()
	params.Pattern = fmt.Sprintf("%s/%s", agt.agentConfig.ArtifactoryRepo, filename)

	resp, _ := agt.artifactoryManager.SearchFiles(params)

	if len(resp) > 0 {
		exists = true
	}
	return exists
}
