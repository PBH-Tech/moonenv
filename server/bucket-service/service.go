package bucketService

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/PBH-Tech/moonenv/handle"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var (
	bucketName = os.Getenv("S3Bucket")
)

type FileData struct {
	B64Str  string `json:"b64String"`
	ObjName string `json:"objectName"`
}

func GetObjectFromS3Bucket(ctx context.Context, req handle.Request, s3Client *s3.Client) (handle.Response, error) {
	status, _ := handle.CheckRequestStatus(ctx)

	if status != nil {
		return *status, nil
	}

	key, found := req.QueryStringParameters["key"]

	if !found {
		return handle.ApiResponse(http.StatusBadRequest, "Query parameter 'key' is required")
	}

	input := &s3.GetObjectInput{
		Bucket: &bucketName,
		Key:    &key,
	}

	result, getErr := s3Client.GetObject(ctx, input)

	if getErr != nil {
		return handle.ApiResponse(http.StatusInternalServerError, "Failed to get object from s3")
	}

	defer result.Body.Close()
	body, err := io.ReadAll(result.Body)

	if err != nil {
		return handle.ApiResponse(http.StatusInternalServerError, "Failed to download object from s3")
	}

	bodyResp := map[string]string{"file": base64.StdEncoding.EncodeToString(body)}

	return handle.ApiResponse(http.StatusOK, bodyResp)
}

func UploadToS3Bucket(ctx context.Context, fileData FileData, s3Client *s3.Client) (handle.Response, error) {
	content, decErr := base64.StdEncoding.DecodeString(fileData.B64Str)

	if decErr != nil {
		return handle.ApiResponse(http.StatusBadRequest, "Invalid base64 string")
	}

	input := &s3.PutObjectInput{Bucket: aws.String(bucketName), Key: aws.String(fileData.ObjName), Body: bytes.NewReader(content)}

	_, putErr := s3Client.PutObject(ctx, input)

	if putErr != nil {
		return handle.ApiResponse(http.StatusInternalServerError, "Failed to upload object to s3")
	}

	respBody := map[string]string{"message": fmt.Sprintf("Object [%v] was uploaded", fileData.ObjName)}

	return handle.ApiResponse(http.StatusOK, respBody)
}
