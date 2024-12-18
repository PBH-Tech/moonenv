package main

import (
	"context"
	"net/http"

	restApi "github.com/PBH-Tech/moonenv/lambdas/util/rest-api"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, req restApi.Request) (restApi.Response, error) {
	switch req.RequestContext.HTTP.Method {
	case http.MethodGet:
		return PullCommand(req)
	case http.MethodPost:
		return PushCommand(req)
	default:
		return restApi.UnhandledMethod()
	}
}
