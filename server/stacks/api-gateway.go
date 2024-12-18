package stacks

import (
	"github.com/aws/aws-cdk-go/awscdk"
	"github.com/aws/aws-cdk-go/awscdk/awsapigatewayv2"
	"github.com/aws/aws-cdk-go/awscdk/awsapigatewayv2integrations"
	"github.com/aws/aws-cdk-go/awscdk/awslambdago"
	"github.com/aws/constructs-go/constructs/v3"
	"github.com/aws/jsii-runtime-go"
)

type CdkApiGatewayProps struct {
	awscdk.StackProps
	CdkLambdaStackFunctions
}

func NewApiGatewayStack(scope constructs.Construct, id string, props *CdkApiGatewayProps) {
	var sProps awscdk.StackProps

	if props != nil {
		sProps = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sProps)

	api := awsapigatewayv2.NewHttpApi(stack, jsii.String("cdk-moonenv-api"), &awsapigatewayv2.HttpApiProps{
		CorsPreflight: &awsapigatewayv2.CorsPreflightOptions{
			AllowOrigins: &[]*string{jsii.String("*")}, //
			AllowMethods: &[]awsapigatewayv2.CorsHttpMethod{
				awsapigatewayv2.CorsHttpMethod_GET,
				awsapigatewayv2.CorsHttpMethod_POST,
			},
		},
	})

	orchestrator := awslambdago.NewGoFunction(stack, jsii.String("orchestrator"), &awslambdago.GoFunctionProps{
		MemorySize: jsii.Number(128),
		Entry:      jsii.String("./lambdas/endpoints/orchestrator"),
		Environment: &map[string]*string{
			"AwsRegion":        props.StackProps.Env.Region,
			"UploadFuncName":   props.CdkLambdaStackFunctions.uploadFileFunc.FunctionArn(),
			"DownloadFuncName": props.CdkLambdaStackFunctions.downloadFileFunc.FunctionArn(),
		},
	})

	api.AddRoutes(&awsapigatewayv2.AddRoutesOptions{
		Path: jsii.String("/orgs/{org}/repos/{repo}"),
		Methods: &[]awsapigatewayv2.HttpMethod{
			awsapigatewayv2.HttpMethod_GET,
			awsapigatewayv2.HttpMethod_POST,
		},
		Integration: awsapigatewayv2integrations.NewHttpLambdaIntegration(jsii.String("orchestrator-get"), orchestrator, &awsapigatewayv2integrations.HttpLambdaIntegrationProps{}),
	})

	props.CdkLambdaStackFunctions.downloadFileFunc.GrantInvoke(orchestrator.Role())
	props.CdkLambdaStackFunctions.uploadFileFunc.GrantInvoke(orchestrator.Role())
}
