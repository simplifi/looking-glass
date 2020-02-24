package downloader

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	sss "github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/simplifi/looking-glass/config"
	"os"
	"path"
)

type s3 struct {
	awsSession session.Session
	awsBucket  string
	awsPrefix  string
}

func newS3(config config.DownloaderConfig) (Downloader, error) {
	err := validateConfig(config)
	if err != nil {
		return nil, err
	}
	awsSess, err := createAwsSession(config.AwsKey, config.AwsSecret, config.AwsRegion)
	if err != nil {
		return nil, err
	}
	downloader := s3{
		awsSession: *awsSess,
		awsBucket:  config.AwsBucket,
		awsPrefix:  config.AwsPrefix,
	}

	return &downloader, nil
}

func validateConfig(config config.DownloaderConfig) error {
	if config.AwsKey == "" {
		return fmt.Errorf("configuration value cannot be empty: AwsKey")
	}
	if config.AwsSecret == "" {
		return fmt.Errorf("configuration value cannot be empty: AwsSecret")
	}
	if config.AwsRegion == "" {
		return fmt.Errorf("configuration value cannot be empty: AwsSecret")
	}
	if config.AwsPrefix == "" {
		return fmt.Errorf("configuration value cannot be empty: AwsKey")
	}
	if config.AwsBucket == "" {
		return fmt.Errorf("configuration value cannot be empty: AwsKey")
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

	resp, err := sss.New(&s3s.awsSession).
		ListObjects(&sss.ListObjectsInput{
			Bucket: aws.String(s3s.awsBucket),
			Prefix: aws.String(s3s.awsPrefix),
		})

	if err != nil {
		return nil, err
	}

	for _, obj := range resp.Contents {
		objects = append(objects, *obj.Key)
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
