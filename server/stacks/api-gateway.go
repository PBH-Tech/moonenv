package stacks

import (
	"github.com/aws/aws-cdk-go/awscdk"
	"github.com/aws/aws-cdk-go/awscdk/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/awslambdago"
	"github.com/aws/aws-cdk-go/awscdk/awsroute53"
	"github.com/aws/aws-cdk-go/awscdk/awsroute53targets"
	"github.com/aws/constructs-go/constructs/v3"
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
		Target:     awsroute53.AddressRecordTarget_FromAlias(awsroute53targets.NewApiGatewayDomain(customDomain)),
	})

	orchestrator := awslambdago.NewGoFunction(stack, jsii.String("MoonenvOrchestrator"), &awslambdago.GoFunctionProps{
		MemorySize:   jsii.Number(128),
		Entry:        jsii.String("./lambdas/endpoints/orchestrator"),
		FunctionName: jsii.String("moonenv-orchestrator"),
		Environment: &map[string]*string{
			"AwsRegion":        props.StackProps.Env.Region,
			"UploadFuncName":   props.CdkLambdaStackFunctions.uploadFileFunc.FunctionArn(),
			"DownloadFuncName": props.CdkLambdaStackFunctions.downloadFileFunc.FunctionArn(),
		},
	})

	// api.AddRoutes(&awsapigatewayv2.AddRoutesOptions{
	// 	Path: jsii.String("/orgs/{org}/repos/{repo}"),
	// 	Methods: &[]awsapigatewayv2.HttpMethod{
	// 		awsapigatewayv2.HttpMethod_GET,
	// 		awsapigatewayv2.HttpMethod_POST,
	// 	},
	// 	Integration: awsapigatewayv2integrations.NewHttpLambdaIntegration(jsii.String("orchestrator"), orchestrator, &awsapigatewayv2integrations.HttpLambdaIntegrationProps{}),
	// })

	createAuthResource(api, props)
	props.CdkLambdaStackFunctions.downloadFileFunc.GrantInvoke(orchestrator.Role())
	props.CdkLambdaStackFunctions.uploadFileFunc.GrantInvoke(orchestrator.Role())

	awscdk.NewCfnOutput(stack, jsii.String("MoonenvApiGatewayUrl"), &awscdk.CfnOutputProps{Value: api.Url()})
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
