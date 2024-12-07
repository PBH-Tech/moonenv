package bucketService

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	restApi "github.com/PBH-Tech/moonenv/lambdas/util/rest-api"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var (
	bucketName = os.Getenv("S3Bucket")
)

type UploadFileData struct {
	B64Str  string
	ObjName string
}

type DownloadFileData struct {
	Key string
}

func GetObjectFromS3Bucket(ctx context.Context, s3Client *s3.Client, key string) (string, error) {
	input := &s3.GetObjectInput{
		Bucket: &bucketName,
		Key:    &key,
	}

	result, getErr := s3Client.GetObject(ctx, input)

	if getErr != nil {
		return "", errors.New("failed to get object from s3")
	}

	defer result.Body.Close()
	body, err := io.ReadAll(result.Body)

	if err != nil {
		return "", errors.New("failed to download object from s3")
	}

	return base64.StdEncoding.EncodeToString(body), nil
}

func UploadToS3Bucket(ctx context.Context, fileData UploadFileData, s3Client *s3.Client) (restApi.Response, error) {
	content, decErr := base64.StdEncoding.DecodeString(fileData.B64Str)

	if decErr != nil {
		return restApi.ApiResponse(http.StatusBadRequest, "Invalid base64 string")
	}

	input := &s3.PutObjectInput{Bucket: aws.String(bucketName), Key: aws.String(fileData.ObjName), Body: bytes.NewReader(content)}

	_, putErr := s3Client.PutObject(ctx, input)

	if putErr != nil {
		return restApi.ApiResponse(http.StatusInternalServerError, "Failed to upload object to s3")
	}

	respBody := map[string]string{"message": fmt.Sprintf("Object [%v] was uploaded", fileData.ObjName)}

	return restApi.ApiResponse(http.StatusOK, respBody)
}
