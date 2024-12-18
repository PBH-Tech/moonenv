package main

import (
	"context"
	"net/http"
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
		code, codeOk   = req.QueryStringParameters["code"]
		state, stateOk = req.QueryStringParameters["state"]
	)

	if !codeOk || !stateOk {
		return restApi.ApiResponse(http.StatusBadRequest, map[string]string{"message": "code and state query parameters are required"})
	}

	return SaveCode(state, code)
}
