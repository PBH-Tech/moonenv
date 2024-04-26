package handle

import (
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
)

var (
	bucketName = os.Getenv("S3Bucket")
)

type Request events.APIGatewayProxyRequest
type Response events.APIGatewayProxyResponse

type ObjKey struct {
	Key string `json:"objectKey"`
}

type FileData struct {
	B64Str  string `json:"b64String"`
	ObjName string `json:"objectName"`
}

func UnhandledMethod() (Response, error) {
	return ApiResponse(http.StatusMethodNotAllowed, "Method not allowed")
}
