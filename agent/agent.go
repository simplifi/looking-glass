package agent

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/jfrog/jfrog-client-go/artifactory"
	"github.com/jfrog/jfrog-client-go/artifactory/auth"
	"github.com/jfrog/jfrog-client-go/artifactory/services"
	aflog "github.com/jfrog/jfrog-client-go/utils/log"
	"github.com/simplifi/looking-glass/config"
	"log"
	"os"
	"path"
	"time"
)

type Agent struct {
	artifactoryManager artifactory.ArtifactoryServicesManager
	awsSession session.Session
	agentConfig config.AgentConfig
	localStoragePath string
}

func New(artifactoryConfig config.ArtifactoryConfig, agentConfig config.AgentConfig) (*Agent, error) {
	artMgr, err := createArtifactoryManager(artifactoryConfig.Url, artifactoryConfig.Key, artifactoryConfig.UserName)
	if err != nil {
		return nil, err
	}
	awsSess, err := createAwsSession(agentConfig.AwsKey, agentConfig.AwsSecret, agentConfig.AwsRegion)
	if err != nil {
		return nil, err
	}
	localStoragePath := path.Join("/tmp", agentConfig.Name)

	agent := Agent{
		artifactoryManager: *artMgr,
		awsSession:         *awsSess,
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

func createAwsSession(awsKey string, awsSecret string, awsRegion string) (*session.Session, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion),
		Credentials: credentials.NewStaticCredentials(awsKey, awsSecret, ""),
	})

	return sess, err
}

func (agt *Agent) Start() {
	for {
		for _, obj := range agt.getS3Objects() {

			// see if any objects are missing from Artifactory
			if !agt.existsInArtifactory(obj) {
				log.Printf("INFO: [mirror] %s", obj)

				// download object to local storage
				localFile := path.Join(agt.localStoragePath, obj)
				agt.downloadS3Obj(obj, localFile)

				// upload to artifactory
				agt.uploadToArtifactory(localFile, obj)

				// clean up temp storage
				err := os.RemoveAll(agt.localStoragePath)
				if err != nil {
					log.Printf("ERROR: Failed to clean up temp storage - %v", err)
				}
			} else {
				log.Printf("INFO: [skip] %s", obj)
			}
		}

		log.Printf("INFO: Sleeping for %d seconds", agt.agentConfig.SleepDuration)
		time.Sleep(time.Duration(agt.agentConfig.SleepDuration)*time.Second)
	}
}

func (agt *Agent) getS3Objects() []string {
	var objects []string

	resp, err := s3.New(&agt.awsSession).
		ListObjects(&s3.ListObjectsInput{
			Bucket: aws.String(agt.agentConfig.AwsBucket),
			Prefix: aws.String(agt.agentConfig.AwsPrefix),
		})

	if err != nil {
		log.Printf("ERROR: Error while listing S3 objects - %v", err)
		return objects
	}

	for _, obj := range resp.Contents {
		objects = append(objects, *obj.Key)
	}

	return objects
}

func (agt *Agent) downloadS3Obj(sourceObj string, targetPath string) {
	downloader := s3manager.NewDownloader(&agt.awsSession)

	// Ensure the temporary download path exists
	err := os.MkdirAll(path.Dir(targetPath), os.ModePerm)
	if err != nil {
		log.Printf("ERROR: failed to create file %q, %v", targetPath, err)
	}

	// Create a file in which we will write the S3 Object contents
	f, err := os.Create(targetPath)
	if err != nil {
		log.Printf("ERROR: failed to create file %q, %v", targetPath, err)
	}

	// Download the object
	_, err = downloader.Download(f, &s3.GetObjectInput{
		Bucket: aws.String(agt.agentConfig.AwsBucket),
		Key:    aws.String(sourceObj),
	})

	if err != nil {
		log.Printf("ERROR: failed to download file, %v", err)
	}
}

func (agt *Agent) uploadToArtifactory(sourceFile string, targetPath string) {
	params := services.NewUploadParams()
	params.Pattern = sourceFile
	params.Target = fmt.Sprintf("%s/%s", agt.agentConfig.ArtifactoryRepo, targetPath)

	_, _, totalFailed, err := agt.artifactoryManager.UploadFiles(params)

	if err != nil || totalFailed > 0 {
		log.Printf("ERROR: failed to upload file %q, %v", sourceFile, err)
	}
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
