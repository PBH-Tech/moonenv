package main

import (
	"context"
	"errors"

	bucketService "github.com/PBH-Tech/moonenv/bucket-service"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var (
	s3Client *s3.Client
)

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, event *bucketService.DownloadFileData) (string, error) {
	cfg, err := config.LoadDefaultConfig(ctx)

	if err != nil {
		return "", errors.New("failed to load SDK Configuration")
	}

	s3Client = s3.NewFromConfig(cfg)

	return bucketService.GetObjectFromS3Bucket(ctx, s3Client, event.Key)
}
