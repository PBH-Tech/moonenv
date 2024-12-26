package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"

	restApi "github.com/PBH-Tech/moonenv/lambdas/util/rest-api"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	CallbackUri              = os.Getenv("CallbackUri")
	CognitoUrl               = os.Getenv("CognitoUrl")
	PollingIntervalInSeconds int64
)

func init() {
	var err error
	PollingIntervalInSeconds, err = strconv.ParseInt(os.Getenv("PollingIntervalInSeconds"), 10, 64)

	if err != nil {
		log.Fatalf("Invalid PollingIntervalInSeconds value: %v", err)

	}

}

func main() {

	lambda.Start(handler)
}

func handler(_ctx context.Context, req restApi.Request) (restApi.Response, error) {
	var (
		clientId                 = req.QueryStringParameters["client_id"]
		deviceCode, deviceCodeOk = req.QueryStringParameters["device_code"]
		grantType, grantTypeOk   = req.QueryStringParameters["grant_type"]
		// Read more about it: https://oauth.net/2/grant-types/device-code/
		deviceCodeGrantType = "urn:ietf:params:oauth:grant-type:device_code"
	)

	/// TODO:improve it
	if !deviceCodeOk {
		return RequestSetOfToken(clientId), nil
	} else if grantTypeOk && deviceCodeOk {
		if grantType == deviceCodeGrantType {
			return RequestJWTs(deviceCode, clientId), nil
		} else {
			return restApi.ApiResponse(http.StatusBadRequest, map[string]string{"message": "Unsupported grant type"}), nil
		}
	} else {
		return restApi.ApiResponse(http.StatusBadRequest, map[string]string{"message": "Missing parameters"}), nil
	}
}
