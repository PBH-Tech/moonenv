package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	bucketService "github.com/PBH-Tech/moonenv/lambdas/util/bucket"
	restApi "github.com/PBH-Tech/moonenv/lambdas/util/rest-api"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	lambdaSdk "github.com/aws/aws-sdk-go/service/lambda"
)

type PushCommandRequest struct {
	B64Str string `json:"b64String"`
}

func PullCommand(req restApi.Request) (restApi.Response, error) {
	pathData := req.PathParameters
	queryDate := req.QueryStringParameters
	pathRequest := bucketService.DownloadFileData{Key: fmt.Sprintf("%s/%s/%s", pathData["org"], pathData["repo"], queryDate["env"])}
	client := getLambdaClient()
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

func getHeader(headers map[string]string, key string) string {
	for hKey, hValue := range headers {
		if strings.EqualFold(hKey, key) {
			return hValue
		}
	}
	return ""
}

func PushCommand(req restApi.Request) (restApi.Response, error) {
	pathData := req.PathParameters
	queryDate := req.QueryStringParameters
	if getHeader(req.Headers, "content-type") != "application/json" {
		return restApi.ApiResponse(http.StatusBadRequest, "Invalid request type")
	}

	var commandData PushCommandRequest

	err := json.Unmarshal([]byte(req.Body), &commandData)

	if err != nil {
		return restApi.ApiResponse(http.StatusBadRequest, "Invalid body request")
	}

	request := bucketService.UploadFileData{B64Str: commandData.B64Str, ObjName: fmt.Sprintf("%s/%s/%s", pathData["org"], pathData["repo"], queryDate["env"])}
	client := getLambdaClient()
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

func getLambdaClient() *lambdaSdk.Lambda {
	newSession := session.Must(session.NewSessionWithOptions(session.Options{SharedConfigState: session.SharedConfigEnable}))

	return lambdaSdk.New(newSession, &aws.Config{Region: aws.String(os.Getenv("AwsRegion"))})
}
