package main

import (
	"context"
	"net/http"

	"github.com/JoaoVitorTeixeira/moonenv-server/handle"
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

func handler(ctx context.Context, req handle.Request) (handle.Response, error) {
	cfg, err := config.LoadDefaultConfig(ctx)

	if err != nil {
		return handle.ApiResponse(http.StatusInternalServerError, "Failed to load SDK Configuration")
	}

	s3Client = s3.NewFromConfig(cfg)

	switch req.HTTPMethod {
	case http.MethodGet:
		return handle.ApiResponse(http.StatusOK, "GET")
	case http.MethodPost:
		return handle.ApiResponse(http.StatusCreated, "POST")
	default:
		return handle.UnhandledMethod()
	}
}
