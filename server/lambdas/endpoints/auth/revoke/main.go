package main

import (
	"context"
	"os"
	"strings"

	restApi "github.com/PBH-Tech/moonenv/lambdas/util/rest-api"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	CognitoUrl = os.Getenv("CognitoUrl")
)

func main() {

	lambda.Start(handler)
}

func handler(_ctx context.Context, req restApi.Request) (restApi.Response, error) {
	var (
		deviceCode = req.QueryStringParameters["device_code"]
		token      = req.Headers["Authorization"]
	)

	return RevokeToken(deviceCode, strings.Replace(token, "Bearer ", "", 1)), nil
}
