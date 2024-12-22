package restApi

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

func ApiResponse(statusCode int, body interface{}) (Response, error) {
	var (
		buf      bytes.Buffer
		respBody []byte
	)

	resp := Response{Headers: map[string]string{"Content-Type": "application/json", "Access-Control-Allow-Origin": "*", "Access-Control-Allow-Methods": "GET, POST"}}
	resp.StatusCode = statusCode

	if body != nil {
		respBody, _ = json.Marshal(body)
		json.HTMLEscape(&buf, respBody)
		resp.Body = string(respBody)
	}

	return resp, nil
}

func UnhandledMethod() (Response, error) {
	return ApiResponse(http.StatusMethodNotAllowed, "Method not allowed")
}

func BuildErrorResponse(statusCode int, message string) (Response, error) {
	return ApiResponse(statusCode, map[string]string{"message": message})
}

type Request events.APIGatewayProxyRequest
type Response events.APIGatewayProxyResponse
