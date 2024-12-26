package main

import (
	"context"
	"os"

	restApi "github.com/PBH-Tech/moonenv/lambdas/util/rest-api"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	StateIndexName = os.Getenv("StateIndexName")
)

func main() {

	lambda.Start(handler)
}

func handler(_ctx context.Context, req restApi.Request) (restApi.Response, error) {
	var (
		code  = req.QueryStringParameters["code"]
		state = req.QueryStringParameters["state"]
	)

	return SaveCode(state, code), nil
}
