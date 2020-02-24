package agent

import (
	"fmt"
	"github.com/simplifi/looking-glass/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAgentNew(t *testing.T) {
	testArtifactoryCfg := config.ArtifactoryConfig{
		URL:      "http://foo.bar",
		UserName: "testing",
		Key:      "123",
	}

	testAgentDownloaderConfig := config.DownloaderConfig{
		Type:      "s3",
		AwsBucket: "test-bucket",
		AwsPrefix: "test-prefix",
		AwsKey:    "MYAWSKEY",
		AwsSecret: "MYAWSSECRET",
		AwsRegion: "us-west-2",
	}

	testAgentConfig := config.AgentConfig{
		Name:            "test",
		ArtifactoryRepo: "test",
		Downloader:      testAgentDownloaderConfig,
		SleepDuration:   100,
	}

	agt, err := New(testArtifactoryCfg, testAgentConfig)
	assert.NoError(t, err)
	assert.NotNil(t, agt)
}

func TestAgentBadDownloadType(t *testing.T) {
	testArtifactoryCfg := config.ArtifactoryConfig{
		URL:      "http://foo.bar",
		UserName: "testing",
		Key:      "123",
	}

	testAgentDownloaderConfig := config.DownloaderConfig{
		Type:      "not-a-valid-type",
		AwsBucket: "test-bucket",
		AwsPrefix: "test-prefix",
		AwsKey:    "MYAWSKEY",
		AwsSecret: "MYAWSSECRET",
		AwsRegion: "us-west-2",
	}

	testAgentConfig := config.AgentConfig{
		Name:            "test",
		ArtifactoryRepo: "test",
		Downloader:      testAgentDownloaderConfig,
		SleepDuration:   100,
	}

	expectedError := fmt.Errorf("unknown type not-a-valid-type")
	_, err := New(testArtifactoryCfg, testAgentConfig)
	if assert.Error(t, err) {
		assert.Equal(t, expectedError, err)
	}
}

func TestAgentMissingRequiredS3Config(t *testing.T) {
	testArtifactoryCfg := config.ArtifactoryConfig{
		URL:      "http://foo.bar",
		UserName: "testing",
		Key:      "123",
	}

	testAgentDownloaderConfig := config.DownloaderConfig{
		Type:      "s3",
		AwsBucket: "test-bucket",
		AwsPrefix: "test-prefix",
		AwsSecret: "MYAWSSECRET",
		AwsRegion: "us-west-2",
	}

	testAgentConfig := config.AgentConfig{
		Name:            "test",
		ArtifactoryRepo: "test",
		Downloader:      testAgentDownloaderConfig,
		SleepDuration:   100,
	}

	expectedError := fmt.Errorf("configuration value cannot be empty: AwsKey")
	_, err := New(testArtifactoryCfg, testAgentConfig)
	if assert.Error(t, err) {
		assert.Equal(t, expectedError, err)
	}
}