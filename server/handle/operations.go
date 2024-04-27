package handle

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

type Request events.APIGatewayProxyRequest
type Response events.APIGatewayProxyResponse

type ObjKey struct {
	Key string `json:"objectKey"`
}

func UnhandledMethod() (Response, error) {
	return ApiResponse(http.StatusMethodNotAllowed, "Method not allowed")
}

func CheckRequestStatus(ctx context.Context) (*Response, error) {
	if ctx.Err() != nil {
		response, _ := ApiResponse(http.StatusRequestTimeout, "Request Canceled")

		return &response, nil
	}

	return nil, nil
}
