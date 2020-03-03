package config

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

func TestConfigRead(t *testing.T) {
	content := []byte(`
---
artifactory:
  url: http://my.artifactory.server/artifactory/
  username: my-artifactory-user
  key: my-artifactory-key
agents:
  - name: my-agent-name
    artifactory_repo: my-repo
    sleep_duration: 900
    downloader:
      type: s3
      config:
        aws_bucket: my-s3-bucket
        aws_key: my-aws-key
        aws_secret: my-aws-secret
        aws_prefix: my-prefix
        aws_region: us-west-2
`)
	tmpfile, _ := ioutil.TempFile("", "config")

	defer os.Remove(tmpfile.Name()) // clean up
	defer tmpfile.Close()
	tmpfile.Write(content)

	cfg, err := Read(tmpfile.Name())
	assert.NoError(t, err)

	assert.Equal(t, "http://my.artifactory.server/artifactory/", cfg.Artifactory.URL)
	assert.Equal(t, "my-artifactory-user", cfg.Artifactory.UserName)
	assert.Equal(t, "my-artifactory-key", cfg.Artifactory.Key)
	assert.Equal(t, "my-agent-name", cfg.Agents[0].Name)
	assert.Equal(t, "s3", cfg.Agents[0].Downloader.Type)
	assert.Equal(t, 900, cfg.Agents[0].SleepDuration)
	assert.Equal(t, "my-repo", cfg.Agents[0].ArtifactoryRepo)

	expectedConfig := map[interface{}]interface{}{
		"aws_bucket": "my-s3-bucket",
		"aws_key":    "my-aws-key",
		"aws_secret": "my-aws-secret",
		"aws_prefix": "my-prefix",
		"aws_region": "us-west-2",
	}
	assert.Equal(t, expectedConfig, cfg.Agents[0].Downloader.Config)
}
