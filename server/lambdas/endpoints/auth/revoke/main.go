package main

import (
	"context"
	"net/http"
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
		deviceCode, deviceCodeOk = req.QueryStringParameters["device_code"]
		token, tokenOk           = req.Headers["Authorization"]
	)

	if !tokenOk {
		return restApi.BuildErrorResponse(http.StatusUnauthorized, "token is missing"), nil
	}

	if !deviceCodeOk {
		return restApi.BuildErrorResponse(http.StatusBadRequest, "device code parameter is required"), nil
	}

	return RevokeToken(deviceCode, strings.Replace(token, "Bearer ", "", 1)), nil
}
