package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	bucketService "github.com/PBH-Tech/moonenv/bucket-service"
	"github.com/PBH-Tech/moonenv/handle"
	orchestratorService "github.com/PBH-Tech/moonenv/orchestrator-service"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	lambdaSdk "github.com/aws/aws-sdk-go/service/lambda"
)

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, req handle.Request) (handle.Response, error) {
	if req.Headers["Content-Type"] != "application/json" {
		return handle.ApiResponse(http.StatusBadRequest, "Invalid request type")
	}

	var commandData orchestratorService.CommandRequest

	err := json.Unmarshal([]byte(req.Body), &commandData)

	if err != nil {
		return handle.ApiResponse(http.StatusBadRequest, "Invalid body request")
	}

	newSession := session.Must(session.NewSessionWithOptions(session.Options{SharedConfigState: session.SharedConfigEnable}))
	client := lambdaSdk.New(newSession, &aws.Config{Region: aws.String(os.Getenv("AwsRegion"))})
	request := bucketService.FileData{B64Str: commandData.B64Str, ObjName: fmt.Sprintf("%s/%s/%s", commandData.Org, commandData.Repo, commandData.Env)}

	payload, err := json.Marshal(request)

	if err != nil {
		return handle.ApiResponse(http.StatusInternalServerError, "Failed while preparing the payload")
	}

	result, err := client.Invoke(&lambdaSdk.InvokeInput{FunctionName: aws.String(os.Getenv("BucketFuncName")), Payload: payload})

	if err != nil {
		return handle.ApiResponse(http.StatusInternalServerError, "Failed invoking function")
	}

	return handle.ApiResponse(int(*result.StatusCode), map[string]string{"Message": "File uploaded"})
}
