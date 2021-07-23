package downloader

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	sss "github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/mitchellh/mapstructure"
	"github.com/simplifi/looking-glass/pkg/looking-glass/config"
)

type s3Config struct {
	AwsBucket string `mapstructure:"aws_bucket"`
	AwsPrefix string `mapstructure:"aws_prefix"`
	AwsKey    string `mapstructure:"aws_key"`
	AwsSecret string `mapstructure:"aws_secret"`
	AwsRegion string `mapstructure:"aws_region"`
}

type s3 struct {
	awsSession session.Session
	awsBucket  string
	awsPrefix  string
}

func newS3(config config.DownloaderConfig) (Downloader, error) {
	var cfg s3Config
	err := mapstructure.Decode(config.Config, &cfg)
	if err != nil {
		return nil, err
	}

	err = validateS3Config(cfg)
	if err != nil {
		return nil, err
	}

	awsSess, err := createAwsSession(cfg.AwsKey, cfg.AwsSecret, cfg.AwsRegion)
	if err != nil {
		return nil, err
	}
	downloader := s3{
		awsSession: *awsSess,
		awsBucket:  cfg.AwsBucket,
		awsPrefix:  cfg.AwsPrefix,
	}

	return &downloader, nil
}

func validateS3Config(cfg s3Config) error {
	requiredConfigs := map[string]string{
		"AwsKey":    cfg.AwsKey,
		"AwsSecret": cfg.AwsSecret,
		"AwsRegion": cfg.AwsRegion,
		"AwsPrefix": cfg.AwsPrefix,
		"AwsBucket": cfg.AwsBucket,
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

func createAwsSession(awsKey string, awsSecret string, awsRegion string) (*session.Session, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(awsRegion),
		Credentials: credentials.NewStaticCredentials(awsKey, awsSecret, ""),
	})

	return sess, err
}

// ListObjects lists the objects available in the S3 bucket
func (s3s *s3) ListObjects() ([]string, error) {
	var objects []string
	path := &sss.ListObjectsV2Input{
		Bucket: aws.String(s3s.awsBucket),
		Prefix: aws.String(s3s.awsPrefix),
	}

	err := sss.New(&s3s.awsSession).
		ListObjectsV2Pages(path, func(page *sss.ListObjectsV2Output, lastPage bool) bool {
			for _, obj := range page.Contents {
				objects = append(objects, *obj.Key)
			}
			return !lastPage
		})

	if err != nil {
		return nil, err
	}

	return objects, nil
}

// GetObject downloads the object specified in sourceObj to the targetPath
func (s3s *s3) GetObject(sourceObj string, targetPath string) error {
	downloader := s3manager.NewDownloader(&s3s.awsSession)

	// Ensure the temporary download path exists
	err := os.MkdirAll(path.Dir(targetPath), os.ModePerm)
	if err != nil {
		return err
	}

	// Create a file in which we will write the S3 Object contents
	f, err := os.Create(targetPath)
	if err != nil {
		return err
	}

	// Download the object
	_, err = downloader.Download(f, &sss.GetObjectInput{
		Bucket: aws.String(s3s.awsBucket),
		Key:    aws.String(sourceObj),
	})
	if err != nil {
		return err
	}

	return nil
}
