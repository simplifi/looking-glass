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
  - name: my-github-release-agent
    artifactory_repo: my-repo
    sleep_duration: 900
    downloader:
      type: github
      config:
        github_repo: simplifi/looking-glass
        github_token: my-github-token
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

// DownloaderConfig holds the configuration for the various downloaders
type DownloaderConfig struct {
	Type   string      `mapstructure:"type"`
	Config interface{} `mapstructure:"config"`
}

// AgentConfig holds Agent specific configuration
type AgentConfig struct {
	Name            string           `mapstructure:"name"`
	ArtifactoryRepo string           `mapstructure:"artifactory_repo"`
	Downloader      DownloaderConfig `mapstructure:"downloader"`
	SleepDuration   int              `mapstructure:"sleep_duration"`
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
