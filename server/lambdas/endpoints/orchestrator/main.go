package main

import (
	"context"
	"net/http"

	"github.com/PBH-Tech/moonenv/lambdas/handle"
	orchestratorService "github.com/PBH-Tech/moonenv/lambdas/orchestrator-service"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, req handle.Request) (handle.Response, error) {
	switch req.RequestContext.HTTP.Method {
	case http.MethodGet:
		return orchestratorService.PullCommand(req)
	case http.MethodPost:
		return orchestratorService.PushCommand(req)
	default:
		return handle.UnhandledMethod()
	}
}
