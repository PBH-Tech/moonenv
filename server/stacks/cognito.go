package stacks

import (
	"github.com/aws/aws-cdk-go/awscdk"
	"github.com/aws/aws-cdk-go/awscdk/awscognito"
	"github.com/aws/constructs-go/constructs/v3"
	"github.com/aws/jsii-runtime-go"
)

type CdkCognitoStackProps struct {
	awscdk.StackProps
}

type CdkCognitoStackResource struct {
}

func NewCognitoStack(scope constructs.Construct, id string, props *CdkCognitoStackProps) (awscognito.UserPool, error) {
	var sProps awscdk.StackProps

	if props != nil {
		sProps = props.StackProps
	}

	stack := awscdk.NewStack(scope, &id, &sProps)

	userPool := awscognito.NewUserPool(stack, jsii.String("moonenv-user-poll"), &awscognito.UserPoolProps{
		UserPoolName:        jsii.String("moonenv-user-pool"),
		SignInCaseSensitive: jsii.Bool(false),
		SignInAliases: &awscognito.SignInAliases{
			Email: jsii.Bool(true),
		},
		AutoVerify: &awscognito.AutoVerifiedAttrs{
			Email: jsii.Bool(true),
		},
	})

	return userPool, nil
}
