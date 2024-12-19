package main

import (
	"context"
	"net/http"
	"os"

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
		deviceCode, deviceCodeOk = req.QueryStringParameters["device_code"]
	)

	if !deviceCodeOk {
		return restApi.BuildErrorResponse(http.StatusBadRequest, "device code parameter is required")
	}

	return RefreshToken(deviceCode, req)
}
