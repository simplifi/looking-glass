package downloader

import (
	"fmt"
	"github.com/simplifi/looking-glass/config"
)

// Downloader downloads objects from various sources
type Downloader interface {
	ListObjects() ([]string, error)
	GetObject(string, string) error
}

// New Downloader, pass in the DownloaderConfig
func New(config config.DownloaderConfig) (Downloader, error) {
	switch config.Type {
	case "s3":
		return newS3(config)
	case "github":
		return newGithub(config)
	default:
		return nil, fmt.Errorf("unknown type %s", config.Type)
	}
}
