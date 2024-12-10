package stacks

import (
	"fmt"

	"github.com/aws/aws-cdk-go/awscdk"
	"github.com/aws/aws-cdk-go/awscdk/awscognito"
	"github.com/aws/constructs-go/constructs/v3"
	"github.com/aws/jsii-runtime-go"
)

type CdkCognitoStackProps struct {
	awscdk.StackProps
}

type CdkCognitoStackResource struct {
	awscognito.CfnUserPoolClient
	awscognito.UserPool
}

var CognitoFullUrl string

func NewCognitoStack(scope constructs.Construct, id string, props *CdkCognitoStackProps) *CdkCognitoStackResource {
	var sProps awscdk.StackProps

	if props != nil {
		sProps = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sProps)

	userPool := awscognito.NewUserPool(stack, jsii.String("MoonenvUserPoll"), &awscognito.UserPoolProps{
		UserPoolName:        jsii.String("moonenv-user-pool"),
		SignInCaseSensitive: jsii.Bool(false),
		SignInAliases: &awscognito.SignInAliases{
			Email: jsii.Bool(true),
		},
		AutoVerify: &awscognito.AutoVerifiedAttrs{
			Email: jsii.Bool(true),
		},
	})
	userPoolId := userPool.UserPoolId()

	cognitoDomain := userPool.AddDomain(jsii.String("MoonenvCognitoDomain"), &awscognito.UserPoolDomainOptions{
		CognitoDomain: &awscognito.CognitoDomainOptions{
			DomainPrefix: jsii.String("moonenv"),
		},
	})
	CognitoFullUrl = fmt.Sprintf("https://%s.auth.%s.amazoncognito.com", *cognitoDomain.DomainName(), *props.Env.Region)

	poolClient := awscognito.NewCfnUserPoolClient(stack, jsii.String("MoonenvPoolClient"), &awscognito.CfnUserPoolClientProps{
		ClientName:                      jsii.String("moonenv-main-client"),
		ExplicitAuthFlows:               jsii.Strings("ALLOW_REFRESH_TOKEN_AUTH", "ALLOW_USER_PASSWORD_AUTH"),
		AllowedOAuthFlows:               jsii.Strings("code"),
		AllowedOAuthScopes:              jsii.Strings("openid", "profile"),
		CallbackUrLs:                    jsii.Strings("http://localhost:3000"),
		AllowedOAuthFlowsUserPoolClient: jsii.Bool(true),
		UserPoolId:                      userPoolId,
		SupportedIdentityProviders:      jsii.Strings("COGNITO"),
	})

	awscdk.NewCfnOutput(stack, jsii.String("MoonenvPoolClientId"), &awscdk.CfnOutputProps{Value: poolClient.GetAtt(jsii.String("ClientId")).ToString()})
	awscdk.NewCfnOutput(stack, jsii.String("MoonenvCognitoDomain"), &awscdk.CfnOutputProps{Value: jsii.String(CognitoFullUrl)})
	awscdk.NewCfnOutput(stack, jsii.String("MoonenvUserPool"), &awscdk.CfnOutputProps{Value: userPoolId})
	return &CdkCognitoStackResource{
		CfnUserPoolClient: poolClient,
		UserPool:          userPool,
	}
}
