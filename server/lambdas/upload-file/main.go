package main

import (
	"context"
	"net/http"

	bucketService "github.com/PBH-Tech/moonenv/lambdas/util/bucket"
	restApi "github.com/PBH-Tech/moonenv/lambdas/util/rest-api"
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

func handler(ctx context.Context, event *bucketService.UploadFileData) restApi.Response {
	cfg, err := config.LoadDefaultConfig(ctx)

	if err != nil {
		return restApi.ApiResponse(http.StatusInternalServerError, "Failed to load SDK Configuration")
	}

	s3Client = s3.NewFromConfig(cfg)

	return bucketService.UploadToS3Bucket(ctx, *event, s3Client)

}
