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
    aws_bucket: my-s3-bucket
    aws_key: my-aws-key
    aws_secret: my-aws-secret
    aws_prefix: my-prefix
    aws_region: us-west-2
    artifactory_repo: my-repo
    sleep_duration: 900
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
	assert.Equal(t, "my-s3-bucket", cfg.Agents[0].AwsBucket)
	assert.Equal(t, "my-prefix", cfg.Agents[0].AwsPrefix)
	assert.Equal(t, "my-aws-key", cfg.Agents[0].AwsKey)
	assert.Equal(t, "my-aws-secret", cfg.Agents[0].AwsSecret)
	assert.Equal(t, "us-west-2", cfg.Agents[0].AwsRegion)
	assert.Equal(t, 900, cfg.Agents[0].SleepDuration)
	assert.Equal(t, "my-repo", cfg.Agents[0].ArtifactoryRepo)

}
