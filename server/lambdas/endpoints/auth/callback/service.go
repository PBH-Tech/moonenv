package main

import (
	"net/http"

	tokenCode "github.com/PBH-Tech/moonenv/lambdas/endpoints/auth"
	restApi "github.com/PBH-Tech/moonenv/lambdas/util/rest-api"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func SaveCode(state string, code string) (restApi.Response, error) {
	tokens, err := tokenCode.QueryToken(StateIndexName, map[string]*dynamodb.Condition{
		"state": {
			ComparisonOperator: aws.String("EQ"),
			AttributeValueList: []*dynamodb.AttributeValue{
				{
					S: aws.String(state),
				},
			},
		},
	})

	if err != nil || tokens == nil || len(tokens) < 1 {
		return restApi.ApiResponse(http.StatusNotFound, "State was not found")
	}

	err = tokenCode.UpdateToken(tokens[0].DeviceCode, tokenCode.TokenCode{LoginCode: code})

	if err != nil {
		return restApi.ApiResponse(http.StatusInternalServerError, "Problem while saving the login code")
	}

	return restApi.ApiResponse(http.StatusNoContent, nil)
}
