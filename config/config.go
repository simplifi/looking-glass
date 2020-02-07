package config

import (
	"github.com/spf13/viper"
	"os"
)

/*
Example configuration file:
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
*/

// Config is used to store configuration for the Agents
type Config struct {
	Artifactory ArtifactoryConfig `mapstructure:"artifactory"`
	Agents      []AgentConfig     `mapstructure:"agents"`
}

// ArtifactoryConfig holds Artifactory specific configuration
type ArtifactoryConfig struct {
	URL      string `mapstructure:"url"`
	UserName string `mapstructure:"username"`
	Key      string `mapstructure:"key"`
}

// AgentConfig holds Agent specific configuration
type AgentConfig struct {
	Name            string `mapstructure:"name"`
	ArtifactoryRepo string `mapstructure:"artifactory_repo"`
	AwsBucket       string `mapstructure:"aws_bucket"`
	AwsPrefix       string `mapstructure:"aws_prefix"`
	AwsKey          string `mapstructure:"aws_key"`
	AwsSecret       string `mapstructure:"aws_secret"`
	AwsRegion       string `mapstructure:"aws_region"`
	SleepDuration   int    `mapstructure:"sleep_duration"`
}

// Read a config file and return a Config
func Read(configPath string) (*Config, error) {
	configFile, readErr := os.Open(configPath)
	if readErr != nil {
		return nil, readErr
	}

	viper.SetConfigType("yaml")
	parseErr := viper.ReadConfig(configFile)
	if parseErr != nil {
		return nil, parseErr
	}

	config := &Config{}

	unmarshalErr := viper.Unmarshal(config)
	if unmarshalErr != nil {
		return nil, unmarshalErr
	}

	return config, nil
}
