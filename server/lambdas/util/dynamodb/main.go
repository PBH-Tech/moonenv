package dynamodb

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	dynamodbService "github.com/aws/aws-sdk-go/service/dynamodb"
)

func NewDynamodb() (*dynamodb.DynamoDB, error) {
	Session, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	})

	if err != nil {
		return nil, err
	}

	return dynamodbService.New(Session), nil
}
