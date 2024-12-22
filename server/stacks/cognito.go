package stacks

import (
	"github.com/aws/aws-cdk-go/awscdk"
	"github.com/aws/aws-cdk-go/awscdk/awscognito"
	"github.com/aws/aws-cdk-go/awscdk/awsroute53"
	"github.com/aws/aws-cdk-go/awscdk/awsroute53targets"
	"github.com/aws/constructs-go/constructs/v3"
	"github.com/aws/jsii-runtime-go"
)

type CdkCognitoStackProps struct {
	awscdk.StackProps
	AuthSubdomain *string
	CdkRoute53StackResource
}

type CdkCognitoStackResource struct {
	awscognito.CfnUserPoolClient
}

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
	/**
	* The parent domain must have a valid DNS A record.
	* Ex.: for auth.moonenv.link, moonenv.link has a A record for 8.8.8.8
	* Read more: https://docs.aws.amazon.com/cognito/latest/developerguide/cognito-user-pools-add-custom-domain.html#cognito-user-pools-add-custom-domain-adding
	 */
	userPoolDomain := userPool.AddDomain(jsii.String("MoonenvCognitoDomain"), &awscognito.UserPoolDomainOptions{
		CustomDomain: &awscognito.CustomDomainOptions{
			DomainName:  props.AuthSubdomain,
			Certificate: props.CdkRoute53StackResource.UsEast1Certificate,
		},
	})
	poolClient := awscognito.NewCfnUserPoolClient(stack, jsii.String("MoonenvPoolClient"), &awscognito.CfnUserPoolClientProps{
		ClientName:                      jsii.String("moonenv-main-client"),
		ExplicitAuthFlows:               jsii.Strings("ALLOW_REFRESH_TOKEN_AUTH", "ALLOW_USER_PASSWORD_AUTH"),
		AllowedOAuthFlows:               jsii.Strings("code"),
		AllowedOAuthScopes:              jsii.Strings("openid", "profile"),
		AllowedOAuthFlowsUserPoolClient: jsii.Bool(true),
		UserPoolId:                      userPoolId,
		SupportedIdentityProviders:      jsii.Strings("COGNITO"),
		EnableTokenRevocation:           jsii.Bool(true),
	})

	awsroute53.NewARecord(stack, jsii.String("MoonenvUserPoolARecord"), &awsroute53.ARecordProps{
		Zone:       props.CdkRoute53StackResource.IHostedZone,
		RecordName: props.AuthSubdomain,
		Target:     awsroute53.AddressRecordTarget_FromAlias(awsroute53targets.NewUserPoolDomainTarget(userPoolDomain)),
	})
	awscdk.NewCfnOutput(stack, jsii.String("MoonenvPoolClientId"), &awscdk.CfnOutputProps{Value: poolClient.GetAtt(jsii.String("ClientId")).ToString()})
	awscdk.NewCfnOutput(stack, jsii.String("MoonenvUserPool"), &awscdk.CfnOutputProps{Value: userPoolId})
	awscdk.NewCfnOutput(stack, jsii.String("MoonenvCognitoDomain"), &awscdk.CfnOutputProps{Value: userPoolDomain.DomainName()})

	return &CdkCognitoStackResource{
		CfnUserPoolClient: poolClient,
	}
}
