package utils

import (
    "context"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/feature/s3/manager"
    "github.com/aws/aws-sdk-go-v2/service/s3"
    "io"
    "log"
)

func UploadFileToS3(bucket, key string, file io.Reader) (string, error) {
    cfg, err := config.LoadDefaultConfig(context.TODO())
    if err != nil {
        log.Printf("Error loading AWS config: %v", err)
        return "", err
    }

    client := s3.NewFromConfig(cfg)
    uploader := manager.NewUploader(client)

    _, err = uploader.Upload(context.TODO(), &s3.PutObjectInput{
        Bucket: &bucket,
        Key:    &key,
        Body:   file,
    })
    if err != nil {
        log.Printf("Error uploading file to S3: %v", err)
        return "", err
    }

    // Construct the URL (replace with your S3 bucket region and URL format)
    url := "https://" + bucket + ".s3.amazonaws.com/" + key
    return url, nil
}