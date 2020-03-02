package downloader

import (
	"fmt"
	"github.com/google/go-github/v29/github"
	"github.com/simplifi/looking-glass/config"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
)

type githubDownloader struct {
	client    github.Client
	repoOwner string
	repoName  string
}

// newGithub returns an initialized githubDownloader struct
func newGithub(config config.DownloaderConfig) (Downloader, error) {
	err := validateGithubConfig(config)
	if err != nil {
		return nil, err
	}

	client := createGithubClient(config.GithubToken)

	// Repo is stored in the configuration as "owner/repo_name" so we split it out here
	repo := strings.Split(config.GithubRepo, "/")

	downloader := githubDownloader{
		client:    *client,
		repoOwner: repo[0],
		repoName:  repo[1],
	}

	return &downloader, nil
}

// validateGithubConfig validates the the configuration is not missing any required values
func validateGithubConfig(config config.DownloaderConfig) error {
	requiredConfigs := map[string]string{
		"GithubRepo": config.GithubRepo,
	}

	var missingConfigs []string

	// Check for configs that are not set
	for cfgName, cfgValue := range requiredConfigs {
		if cfgValue == "" {
			missingConfigs = append(missingConfigs, cfgName)
		}
	}

	// Error on all the missing configs
	if len(missingConfigs) > 0 {
		return fmt.Errorf("configuration values cannot be empty: %s", strings.Join(missingConfigs, ", "))
	}

	return nil
}

// createGithubClient creates a new Github client to be used by the github downloader
func createGithubClient(accessToken string) *github.Client {
	ctx := context.Background()
	if accessToken != "" {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: accessToken},
		)
		tc := oauth2.NewClient(ctx, ts)
		return github.NewClient(tc)
	}

	return github.NewClient(nil)
}

// buildObjectPath builds a path to an object that matches the Artifactory path
func (ghd *githubDownloader) buildObjectPath(tagName string, assetName string) string {
	return fmt.Sprintf("%s/%s/%s/%s", ghd.repoOwner, ghd.repoName, tagName, assetName)
}

// ListObjects lists the objects available in the Github Repo
func (ghd *githubDownloader) ListObjects() ([]string, error) {
	var objects []string

	ctx := context.Background()

	releases, _, err := ghd.client.Repositories.ListReleases(ctx, ghd.repoOwner, ghd.repoName, nil)
	if err != nil {
		return nil, err
	}

	for _, release := range releases {
		assets, _, err := ghd.client.Repositories.ListReleaseAssets(ctx, ghd.repoOwner, ghd.repoName, *release.ID, nil)
		if err != nil {
			return nil, err
		}

		for _, asset := range assets {
			objects = append(objects, ghd.buildObjectPath(*release.TagName, *asset.Name))
		}
	}
	return objects, nil
}

// getReleaseID Gets the ID of a release from the release tag
func (ghd *githubDownloader) getReleaseID(releaseTag string) (int64, error) {
	ctx := context.Background()

	releases, _, err := ghd.client.Repositories.ListReleases(ctx, ghd.repoOwner, ghd.repoName, nil)
	if err != nil {
		return -1, err
	}

	for _, release := range releases {
		if *release.TagName == releaseTag {
			return *release.ID, nil
		}
	}
	return -1, fmt.Errorf("release '%s' not found", releaseTag)
}

// getAssetID Gets the ID of a asset from the provided assetName for the given releaseID
func (ghd *githubDownloader) getAssetID(releaseID int64, assetName string) (int64, error) {
	ctx := context.Background()

	assets, _, err := ghd.client.Repositories.ListReleaseAssets(ctx, ghd.repoOwner, ghd.repoName, releaseID, nil)
	if err != nil {
		return -1, err
	}

	for _, asset := range assets {
		if *asset.Name == assetName {
			return *asset.ID, nil
		}
	}
	return -1, fmt.Errorf("asset '%s' not found", assetName)
}

// GetObject downloads the object specified in sourceObj to the targetPath
func (ghd *githubDownloader) GetObject(sourceObj string, targetPath string) error {
	// Ensure the temporary download path exists
	err := os.MkdirAll(path.Dir(targetPath), os.ModePerm)
	if err != nil {
		return err
	}

	// Create a file in which we will write the github release
	f, err := os.Create(targetPath)
	if err != nil {
		return err
	}

	// Identify the ReleaseID and AssetID
	ctx := context.Background()
	pathSplit := strings.Split(sourceObj, "/")

	releaseTag := pathSplit[2]
	assetName := pathSplit[3]

	releaseID, err := ghd.getReleaseID(releaseTag)
	if err != nil {
		return err
	}

	assetID, err := ghd.getAssetID(releaseID, assetName)
	if err != nil {
		return err
	}

	// Download the asset
	rc, _, err := ghd.client.Repositories.DownloadReleaseAsset(ctx, ghd.repoOwner, ghd.repoName, assetID, http.DefaultClient)
	_, err = io.Copy(f, rc)

	if err != nil {
		return err
	}

	return nil
}
