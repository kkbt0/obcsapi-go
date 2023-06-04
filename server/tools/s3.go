package tools

import (
	"bytes"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func S3FileStore(file_key string, inFileBytes []byte) (string, error) {
	sess, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(NowRunConfig.S3Compatible.AccessKey, NowRunConfig.S3Compatible.SecretKey, ""),
		Endpoint:    aws.String(NowRunConfig.S3Compatible.EndPoint),
		Region:      aws.String(NowRunConfig.S3Compatible.Region),
	})
	if err != nil {
		log.Println(err)
		return "", err
	}
	svc := s3.New(sess)

	// Specify the parameters for the S3 PutObject API
	params := &s3.PutObjectInput{
		Bucket: aws.String(NowRunConfig.S3Compatible.Bucket),
		Key:    aws.String(file_key),
		Body:   bytes.NewReader(inFileBytes),
	}
	// Upload the file to S3
	_, err = svc.PutObject(params)
	if err != nil {
		log.Println(err)
		return "", err
	}
	fmt.Printf("https://%s/%s/%s\n", NowRunConfig.S3Compatible.EndPoint, NowRunConfig.S3Compatible.Bucket, file_key)
	if NowRunConfig.S3Compatible.BaseUrl == "" {
		return fmt.Sprintf("https://%s/%s/%s", NowRunConfig.S3Compatible.EndPoint, NowRunConfig.S3Compatible.Bucket, file_key), nil
	} else {
		return fmt.Sprintf("%s/%s", NowRunConfig.S3Compatible.BaseUrl, file_key), nil
	}
}
