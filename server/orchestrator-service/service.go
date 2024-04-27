package orchestratorService

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	bucketService "github.com/PBH-Tech/moonenv/bucket-service"
	"github.com/PBH-Tech/moonenv/handle"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	lambdaSdk "github.com/aws/aws-sdk-go/service/lambda"
)

type PushCommandRequest struct {
	Org    string `json:"org"`
	Repo   string `json:"repo"`
	Env    string `json:"env"`
	B64Str string `json:"b64String"`
}

func PullCommand(req handle.Request) (handle.Response, error) {

	commandData := req.QueryStringParameters
	request := bucketService.DownloadFileData{Key: fmt.Sprintf("%s/%s/%s", commandData["org"], commandData["repo"], commandData["env"])}
	client := getLambdaClient()
	payload, err := json.Marshal(request)

	if err != nil {
		return handle.ApiResponse(http.StatusInternalServerError, "Failed while preparing the payload")
	}

	result, err := client.Invoke(&lambdaSdk.InvokeInput{Payload: payload, FunctionName: aws.String(os.Getenv("DownloadFuncName"))})

	if err != nil {
		return handle.ApiResponse(http.StatusInternalServerError, "Failed invoking function")
	}

	var response string

	if err := json.Unmarshal(result.Payload, &response); err != nil {
		return handle.ApiResponse(http.StatusInternalServerError, "Failed converting JSON")
	}

	return handle.ApiResponse(http.StatusOK, map[string]string{"file": response})
}

func PushCommand(req handle.Request) (handle.Response, error) {
	if req.Headers["Content-Type"] != "application/json" {
		return handle.ApiResponse(http.StatusBadRequest, "Invalid request type")
	}

	var commandData PushCommandRequest

	err := json.Unmarshal([]byte(req.Body), &commandData)

	if err != nil {
		return handle.ApiResponse(http.StatusBadRequest, "Invalid body request")
	}

	request := bucketService.UploadFileData{B64Str: commandData.B64Str, ObjName: fmt.Sprintf("%s/%s/%s", commandData.Org, commandData.Repo, commandData.Env)}
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
