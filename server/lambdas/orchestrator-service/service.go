package orchestratorService

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	bucketService "github.com/PBH-Tech/moonenv/lambdas/bucket-service"
	"github.com/PBH-Tech/moonenv/lambdas/handle"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	lambdaSdk "github.com/aws/aws-sdk-go/service/lambda"
)

type PushCommandRequest struct {
	B64Str string `json:"b64String"`
}

func PullCommand(req handle.Request) (handle.Response, error) {
	pathData := req.PathParameters
	queryDate := req.QueryStringParameters
	pathRequest := bucketService.DownloadFileData{Key: fmt.Sprintf("%s/%s/%s", pathData["org"], pathData["repo"], queryDate["env"])}
	client := getLambdaClient()
	payload, err := json.Marshal(pathRequest)

	if err != nil {
		return handle.ApiResponse(http.StatusInternalServerError, "Failed while preparing the payload")
	}

	result, err := client.Invoke(&lambdaSdk.InvokeInput{Payload: payload, FunctionName: aws.String(os.Getenv("DownloadFuncName"))})

	if err != nil {
		return handle.ApiResponse(http.StatusInternalServerError, "Failed invoking function")
	}

	var response string

	if err := json.Unmarshal(result.Payload, &response); err != nil {
		return handle.ApiResponse(http.StatusNotFound, "File does not exist")
	}

	return handle.ApiResponse(http.StatusOK, map[string]string{"file": response})
}

func getHeader(headers map[string]string, key string) string {
	for hKey, hValue := range headers {
		if strings.EqualFold(hKey, key) {
			return hValue
		}
	}
	return ""
}

func PushCommand(req handle.Request) (handle.Response, error) {
	pathData := req.PathParameters
	queryDate := req.QueryStringParameters
	if getHeader(req.Headers, "content-type") != "application/json" {
		return handle.ApiResponse(http.StatusBadRequest, "Invalid request type")
	}

	var commandData PushCommandRequest

	err := json.Unmarshal([]byte(req.Body), &commandData)

	if err != nil {
		return handle.ApiResponse(http.StatusBadRequest, "Invalid body request")
	}

	request := bucketService.UploadFileData{B64Str: commandData.B64Str, ObjName: fmt.Sprintf("%s/%s/%s", pathData["org"], pathData["repo"], queryDate["env"])}
	client := getLambdaClient()
	payload, err := json.Marshal(request)

	if err != nil {
		return handle.ApiResponse(http.StatusInternalServerError, "Failed while preparing the payload")
	}

	result, err := client.Invoke(&lambdaSdk.InvokeInput{FunctionName: aws.String(os.Getenv("UploadFuncName")), Payload: payload})

	if err != nil {
		return handle.ApiResponse(http.StatusInternalServerError, "Failed invoking function")
	}

	return handle.ApiResponse(int(*result.StatusCode), map[string]string{"message": "File uploaded"})
}

func getLambdaClient() *lambdaSdk.Lambda {
	newSession := session.Must(session.NewSessionWithOptions(session.Options{SharedConfigState: session.SharedConfigEnable}))

	return lambdaSdk.New(newSession, &aws.Config{Region: aws.String(os.Getenv("AwsRegion"))})
}
