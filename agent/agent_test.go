package agent

import (
	"github.com/simplifi/looking-glass/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAgentNew(t *testing.T) {
	testArtifactoryCfg := config.ArtifactoryConfig{
		Url: "http://foo.bar",
		UserName: "testing",
		Key: "123",
	}

	testAgentConfig := config.AgentConfig{
		Name:            "test",
		ArtifactoryRepo: "test",
		AwsBucket:       "test-bucket",
		AwsPrefix:       "test-prefix",
		AwsKey:          "MYAWSKEY",
		AwsSecret:       "MYAWSSECRET",
		AwsRegion:       "us-west-2",
		SleepDuration:   100,
	}

	agt, err := New(testArtifactoryCfg, testAgentConfig)
	assert.NoError(t, err)
	assert.NotNil(t, agt)
}