package stacks

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsroute53"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsroute53targets"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type CdkApiGatewayProps struct {
	awscdk.StackProps
	CdkLambdaStackFunctions
	CdkRoute53StackResource
	TokenCodeTable          awsdynamodb.Table
	CognitoStack            CdkCognitoStackResource
	TokenCodeStateIndexName *string
	RestApiSubdomain        *string
}

func NewApiGatewayStack(scope constructs.Construct, id string, props *CdkApiGatewayProps) {
	var sProps awscdk.StackProps

	if props != nil {
		sProps = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sProps)

	api := awsapigateway.NewRestApi(stack, jsii.String("MoonenvRestApi"), &awsapigateway.RestApiProps{
		RestApiName: jsii.Sprintf("moonenv-rest-api"),
	})

	customDomain := api.AddDomainName(jsii.String("MoonenvRestApiDomainName"), &awsapigateway.DomainNameOptions{
		Certificate: props.CdkRoute53StackResource.Certificate,
		DomainName:  props.RestApiSubdomain,
	})

	awsroute53.NewARecord(stack, jsii.String("MoonenvRestApiARecord"), &awsroute53.ARecordProps{
		Zone:       props.CdkRoute53StackResource.IHostedZone,
		RecordName: props.RestApiSubdomain,
		Target:     awsroute53.RecordTarget_FromAlias(awsroute53targets.NewApiGatewayDomain(customDomain)),
	})

	createAuthResource(api, props)
	createOrgResource(api, props)

	awscdk.NewCfnOutput(stack, jsii.String("MoonenvApiGatewayUrl"), &awscdk.CfnOutputProps{Value: api.Url()})
}

func createOrgResource(api awsapigateway.RestApi, props *CdkApiGatewayProps) {
	lambdas := props.CdkLambdaStackFunctions
	orgResource := api.Root().AddResource(jsii.String("orgs"), &awsapigateway.ResourceOptions{})
	orgIdResource := orgResource.AddResource(jsii.String("{orgId}"), &awsapigateway.ResourceOptions{})
	repoResource := orgIdResource.AddResource(jsii.String("repos"), &awsapigateway.ResourceOptions{})
	repoIdResource := repoResource.AddResource(jsii.String("{repoId}"), &awsapigateway.ResourceOptions{})

	repoIdResource.AddMethod(jsii.String(*jsii.String("GET")),
		awsapigateway.NewLambdaIntegration(lambdas.pullCommand, &awsapigateway.LambdaIntegrationOptions{}),
		&awsapigateway.MethodOptions{})
	repoIdResource.AddMethod(jsii.String(*jsii.String("POST")),
		awsapigateway.NewLambdaIntegration(lambdas.pushCommand, &awsapigateway.LambdaIntegrationOptions{}),
		&awsapigateway.MethodOptions{})

}

func createAuthResource(api awsapigateway.RestApi, props *CdkApiGatewayProps) {
	callbackUri := GetApiGatewayCallbackUri(props.RestApiSubdomain)
	lambdas := props.CdkLambdaStackFunctions
	authResource := api.Root().AddResource(jsii.String("auth"), &awsapigateway.ResourceOptions{})

	authResource.AddResource(jsii.String("token"), &awsapigateway.ResourceOptions{}).
		AddMethod(jsii.String(*jsii.String("GET")),
			awsapigateway.NewLambdaIntegration(lambdas.tokenAuth, &awsapigateway.LambdaIntegrationOptions{}),
			&awsapigateway.MethodOptions{})

	authResource.AddResource(jsii.String("callback"), &awsapigateway.ResourceOptions{}).
		AddMethod(jsii.String("GET"),
			awsapigateway.NewLambdaIntegration(lambdas.callbackAuth, &awsapigateway.LambdaIntegrationOptions{}), &awsapigateway.MethodOptions{})

	authResource.AddResource(jsii.String("refresh-token"), &awsapigateway.ResourceOptions{}).
		AddMethod(jsii.String("POST"),
			awsapigateway.NewLambdaIntegration(lambdas.refreshTokenAuth, &awsapigateway.LambdaIntegrationOptions{}), &awsapigateway.MethodOptions{})

	authResource.AddResource(jsii.String("revoke"), &awsapigateway.ResourceOptions{}).
		AddMethod(jsii.String("POST"),
			awsapigateway.NewLambdaIntegration(lambdas.revokeTokenAuth, &awsapigateway.LambdaIntegrationOptions{}), &awsapigateway.MethodOptions{})

	props.CognitoStack.SetCallbackUrLs(&[]*string{callbackUri})

}

func GetApiGatewayCallbackUri(restApiSubdomain *string) *string {
	return jsii.Sprintf("https://%s/auth/callback", *restApiSubdomain)
}
