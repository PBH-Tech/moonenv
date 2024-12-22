package main

import (
	"context"

	restApi "github.com/PBH-Tech/moonenv/lambdas/util/rest-api"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handler)
}

func handler(_ctx context.Context, req restApi.Request) (restApi.Response, error) {
	return PushCommand(req)
}
