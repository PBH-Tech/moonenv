package tokenCode

import (
	"os"
	"reflect"

	"github.com/PBH-Tech/moonenv/lambdas/util/dynamodb"
	"github.com/aws/aws-sdk-go-v2/aws"
	dynamodbService "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type TokenCode struct {
	DeviceCode              string `json:"deviceCode"`
	UserCode                string `json:"userCode"`
	AuthorizationUri        string `json:"authorizationUri"`
	VerificationUriComplete string `json:"verificationUriComplete"`
	ClientId                string `json:"clientId"`
	ExpireAt                string `json:"expireAt"`
	LastCheckedAt           string `json:"lastCheckedAt"`
	CodeChallenge           string `json:"code_challenge"`
	CodeVerifier            string `json:"-"` // Omitting it
	// TODO: find a way to turn it into something like an enum
	Status string `json:"status"`
}

var (
	tokenCodeTableName = aws.String(os.Getenv("TokenCodeTableName"))
)

func InsertToken(token TokenCode) (*TokenCode, error) {
	item, err := dynamodbattribute.MarshalMap(token)

	if err != nil {
		return nil, err
	}

	input := &dynamodbService.PutItemInput{
		Item:      item,
		TableName: tokenCodeTableName,
	}

	client, err := dynamodb.NewDynamodb()

	if err != nil {
		return nil, err
	}

	client.PutItem(input)

	return &token, nil
}

func GetToken(deviceCode string) (*TokenCode, error) {
	client, err := dynamodb.NewDynamodb()

	if err != nil {
		return nil, err
	}

	input := &dynamodbService.GetItemInput{
		Key: map[string]*dynamodbService.AttributeValue{
			"deviceCode": {
				S: aws.String(deviceCode),
			},
		},
		TableName: tokenCodeTableName,
	}

	result, err := client.GetItem(input)

	if err != nil {
		return nil, err
	}

	tokenCode := new(TokenCode)

	err = dynamodbattribute.UnmarshalMap(result.Item, tokenCode)

	if result.Item == nil || err != nil {
		return nil, err
	}

	return tokenCode, nil
}

func UpdateToken(deviceCode string, tokenCodeToUpdate TokenCode) error {
	client, err := dynamodb.NewDynamodb()

	if err != nil {
		return err
	}

	item, err := dynamodbattribute.MarshalMap(tokenCodeToUpdate)

	if err != nil {
		return err
	}

	var updateExpression string
	expressionAttributeValues := make(map[string]*dynamodbService.AttributeValue)
	expressionAttributeNames := make(map[string]*string)

	v := reflect.ValueOf(tokenCodeToUpdate)
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		typeField := v.Type().Field(i)
		tag := typeField.Tag.Get("json")

		if tag == "" {
			tag = typeField.Name
		}

		if !field.IsZero() { // Removing fields that is not set, avoiding setting null in the database
			attributePlaceholder := "#" + tag
			valuePlaceholder := ":" + tag

			// Add reserved keyword handling if necessary
			expressionAttributeNames[attributePlaceholder] = aws.String(tag)

			// Construct update expression
			if updateExpression != "" {
				updateExpression += ", "
			}
			updateExpression += attributePlaceholder + " = " + valuePlaceholder

			// Add value
			expressionAttributeValues[valuePlaceholder] = item[tag]
		}
	}

	println("updateExpression: %+v", updateExpression)

	input := dynamodbService.UpdateItemInput{
		Key: map[string]*dynamodbService.AttributeValue{
			"deviceCode": {
				S: aws.String(deviceCode),
			},
		},
		UpdateExpression:          aws.String("SET " + updateExpression),
		ExpressionAttributeValues: expressionAttributeValues,
		ExpressionAttributeNames:  expressionAttributeNames,
		TableName:                 tokenCodeTableName,
	}

	_, err = client.UpdateItem(&input)

	if err != nil {
		println("%s", err.Error())
		return err
	}

	return nil
}
