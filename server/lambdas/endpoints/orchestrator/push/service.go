package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/PBH-Tech/moonenv/lambdas/endpoints/orchestrator"
	bucketService "github.com/PBH-Tech/moonenv/lambdas/util/bucket"
	restApi "github.com/PBH-Tech/moonenv/lambdas/util/rest-api"
	"github.com/aws/aws-sdk-go/aws"
	lambdaSdk "github.com/aws/aws-sdk-go/service/lambda"
)

type PushCommandRequest struct {
	B64Str string `json:"b64String"`
}

func PushCommand(req restApi.Request) restApi.Response {
	pathData := req.PathParameters
	queryDate := req.QueryStringParameters
	if orchestrator.GetHeader(req.Headers, "content-type") != "application/json" {
		return restApi.ApiResponse(http.StatusBadRequest, "Invalid request type")
	}

	var commandData PushCommandRequest

	err := json.Unmarshal([]byte(req.Body), &commandData)

	if err != nil {
		return restApi.ApiResponse(http.StatusBadRequest, "Invalid body request")
	}

	request := bucketService.UploadFileData{B64Str: commandData.B64Str, ObjName: fmt.Sprintf("%s/%s/%s", pathData["orgId"], pathData["repoId"], queryDate["env"])}
	client := orchestrator.GetLambdaClient()
	payload, err := json.Marshal(request)

	if err != nil {
		return restApi.ApiResponse(http.StatusInternalServerError, "Failed while preparing the payload")
	}

	result, err := client.Invoke(&lambdaSdk.InvokeInput{FunctionName: aws.String(os.Getenv("UploadFuncName")), Payload: payload})

	if err != nil {
		return restApi.ApiResponse(http.StatusInternalServerError, "Failed invoking function")
	}

	return restApi.ApiResponse(int(*result.StatusCode), map[string]string{"message": "File uploaded"})
}
