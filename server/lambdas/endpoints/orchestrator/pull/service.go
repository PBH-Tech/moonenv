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

func PullCommand(req restApi.Request) restApi.Response {
	pathData := req.PathParameters
	queryDate := req.QueryStringParameters
	pathRequest := bucketService.DownloadFileData{Key: fmt.Sprintf("%s/%s/%s", pathData["org"], pathData["repo"], queryDate["env"])}
	client := orchestrator.GetLambdaClient()
	payload, err := json.Marshal(pathRequest)

	if err != nil {
		return restApi.ApiResponse(http.StatusInternalServerError, "Failed while preparing the payload")
	}

	result, err := client.Invoke(&lambdaSdk.InvokeInput{Payload: payload, FunctionName: aws.String(os.Getenv("DownloadFuncName"))})

	if err != nil {
		return restApi.ApiResponse(http.StatusInternalServerError, "Failed invoking function")
	}

	var response string

	if err := json.Unmarshal(result.Payload, &response); err != nil {
		return restApi.ApiResponse(http.StatusNotFound, "File does not exist")
	}

	return restApi.ApiResponse(http.StatusOK, map[string]string{"file": response})
}
