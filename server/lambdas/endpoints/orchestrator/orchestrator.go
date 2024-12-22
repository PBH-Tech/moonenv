package orchestrator

import (
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	lambdaSdk "github.com/aws/aws-sdk-go/service/lambda"
)

func GetHeader(headers map[string]string, key string) string {
	for hKey, hValue := range headers {
		if strings.EqualFold(hKey, key) {
			return hValue
		}
	}
	return ""
}

func GetLambdaClient() *lambdaSdk.Lambda {
	newSession := session.Must(session.NewSessionWithOptions(session.Options{SharedConfigState: session.SharedConfigEnable}))

	return lambdaSdk.New(newSession, &aws.Config{Region: aws.String(os.Getenv("AwsRegion"))})
}
