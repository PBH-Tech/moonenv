package handle

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
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

func UploadToS3Bucket(ctx context.Context, req Request, s3Client *s3.Client) (Response, error) {

	if ctx.Err() != nil {
		return ApiResponse(http.StatusRequestTimeout, "Request Canceled")
	}

	var fileData FileData

	if req.Headers["Content-Type"] != "application/json" { // TODO: Find another place to validate it
		return ApiResponse(http.StatusBadRequest, "Invalid request type")
	}

	err := json.Unmarshal([]byte(req.Body), &fileData)

	if err != nil {
		return ApiResponse(http.StatusBadRequest, "Invalid request body")
	}

	content, decErr := base64.StdEncoding.DecodeString(fileData.B64Str)

	if decErr != nil {
		return ApiResponse(http.StatusBadRequest, "Invalid base64 string")
	}

	input := &s3.PutObjectInput{Bucket: aws.String(bucketName), Key: aws.String(fileData.ObjName), Body: bytes.NewReader(content)}

	_, putErr := s3Client.PutObject(ctx, input)

	if putErr != nil {
		return ApiResponse(http.StatusInternalServerError, "Failed to upload object to s3")
	}

	respBody := map[string]string{"message": fmt.Sprintf("Object [%v] was uploaded", fileData.ObjName)}

	return ApiResponse(http.StatusOK, respBody)
}
