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
		Type: "s3",
		Config: map[interface{}]interface{}{
			"aws_bucket": "test-bucket",
			"aws_key":    "MYAWSKEY",
			"aws_prefix": "test-prefix",
			"aws_secret": "MYAWSSECRET",
			"aws_region": "us-west-2",
		},
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
		Type: "not-a-valid-type",
		Config: map[interface{}]interface{}{
			"aws_bucket": "test-bucket",
			"aws_prefix": "test-prefix",
			"aws_secret": "MYAWSSECRET",
			"aws_region": "us-west-2",
		},
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

func TestAgentMissingRequiredS3Configs(t *testing.T) {
	testArtifactoryCfg := config.ArtifactoryConfig{
		URL:      "http://foo.bar",
		UserName: "testing",
		Key:      "123",
	}

	testAgentDownloaderConfig := config.DownloaderConfig{
		Type: "s3",
		Config: map[interface{}]interface{}{
			"aws_bucket": "test-bucket",
			"aws_prefix": "test-prefix",
			"aws_secret": "MYAWSSECRET",
			"aws_region": "us-west-2",
		},
	}

	testAgentConfig := config.AgentConfig{
		Name:            "test",
		ArtifactoryRepo: "test",
		Downloader:      testAgentDownloaderConfig,
		SleepDuration:   100,
	}

	expectedError := fmt.Errorf("configuration values cannot be empty: AwsKey")
	_, err := New(testArtifactoryCfg, testAgentConfig)
	if assert.Error(t, err) {
		assert.Equal(t, expectedError, err)
	}
}
