package stacks

import (
	"fmt"
	"strconv"

	"github.com/aws/aws-cdk-go/awscdk"
	"github.com/aws/aws-cdk-go/awscdk/awsapigatewayv2"
	"github.com/aws/aws-cdk-go/awscdk/awsapigatewayv2integrations"
	"github.com/aws/aws-cdk-go/awscdk/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/awslambdago"
	"github.com/aws/constructs-go/constructs/v3"
	"github.com/aws/jsii-runtime-go"
)

type CdkApiGatewayProps struct {
	awscdk.StackProps
	CdkLambdaStackFunctions
	TokenCodeTable          awsdynamodb.Table
	CognitoStack            CdkCognitoStackResource
	TokenCodeStateIndexName *string
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

	callbackUri := fmt.Sprintf("%sauth/callback", *api.Url())
	//TODO: improve this
	// Tried to get this value from cognito using SSM, but it makes cycle dependency
	cognitoUrl := "https://moonenv.auth.ap-southeast-2.amazoncognito.com"
	tokenAuth := awslambdago.NewGoFunction(stack, jsii.String("MoonenvAuthToken"), &awslambdago.GoFunctionProps{
		MemorySize:   jsii.Number(128),
		Entry:        jsii.String("./lambdas/endpoints/auth/token"),
		FunctionName: jsii.String("moonenv-auth-token"),
		Environment: &map[string]*string{
			"TokenCodeTableName":       props.TokenCodeTable.TableName(),
			"CallbackUri":              jsii.String(callbackUri),
			"PollingIntervalInSeconds": jsii.String(strconv.FormatInt(int64(3), 10)),

			"CognitoUrl": jsii.String(cognitoUrl),
		},
	})
	callbackAuth := awslambdago.NewGoFunction(stack, jsii.Sprintf("MoonenvAuthCallback"), &awslambdago.GoFunctionProps{
		MemorySize:   jsii.Number(128),
		Entry:        jsii.Sprintf("./lambdas/endpoints/auth/callback"),
		FunctionName: jsii.Sprintf("moonenv-auth-callback"),
		Environment: &map[string]*string{
			"StateIndexName":     props.TokenCodeStateIndexName,
			"TokenCodeTableName": props.TokenCodeTable.TableName(),
		},
	})
	refreshTokenAuth := awslambdago.NewGoFunction(stack, jsii.Sprintf("MoonenvAuthRefreshToken"), &awslambdago.GoFunctionProps{
		MemorySize:   jsii.Number(128),
		Entry:        jsii.Sprintf("./lambdas/endpoints/auth/refresh"),
		FunctionName: jsii.Sprintf("moonenv-auth-refresh-token"),
		Environment: &map[string]*string{
			"CognitoUrl":         jsii.String(cognitoUrl),
			"TokenCodeTableName": props.TokenCodeTable.TableName(),
		},
	})
	revokeTokenAuth := awslambdago.NewGoFunction(stack, jsii.Sprintf("MoonenvAuthRevokeToken"), &awslambdago.GoFunctionProps{
		MemorySize:   jsii.Number(128),
		Entry:        jsii.Sprintf("./lambdas/endpoints/auth/revoke"),
		FunctionName: jsii.Sprintf("moonenv-auth-revoke-token"),
		Environment: &map[string]*string{
			"CognitoUrl":         jsii.String(cognitoUrl),
			"TokenCodeTableName": props.TokenCodeTable.TableName(),
		},
	})

	api.AddRoutes(&awsapigatewayv2.AddRoutesOptions{
		Path: jsii.String("/orgs/{org}/repos/{repo}"),
		Methods: &[]awsapigatewayv2.HttpMethod{
			awsapigatewayv2.HttpMethod_GET,
			awsapigatewayv2.HttpMethod_POST,
		},
		Integration: awsapigatewayv2integrations.NewHttpLambdaIntegration(jsii.String("orchestrator"), orchestrator, &awsapigatewayv2integrations.HttpLambdaIntegrationProps{}),
	})
	api.AddRoutes(&awsapigatewayv2.AddRoutesOptions{
		Path:        jsii.Sprintf("/auth/token"),
		Methods:     &[]awsapigatewayv2.HttpMethod{awsapigatewayv2.HttpMethod_GET},
		Integration: awsapigatewayv2integrations.NewHttpLambdaIntegration(jsii.String("auth"), tokenAuth, &awsapigatewayv2integrations.HttpLambdaIntegrationProps{}),
	})
	api.AddRoutes(&awsapigatewayv2.AddRoutesOptions{
		Path:        jsii.Sprintf("/auth/callback"),
		Methods:     &[]awsapigatewayv2.HttpMethod{awsapigatewayv2.HttpMethod_GET},
		Integration: awsapigatewayv2integrations.NewHttpLambdaIntegration(jsii.Sprintf("callback"), callbackAuth, &awsapigatewayv2integrations.HttpLambdaIntegrationProps{}),
	})
	api.AddRoutes(&awsapigatewayv2.AddRoutesOptions{
		Path:        jsii.Sprintf("/auth/refresh-token"),
		Methods:     &[]awsapigatewayv2.HttpMethod{awsapigatewayv2.HttpMethod_POST},
		Integration: awsapigatewayv2integrations.NewHttpLambdaIntegration(jsii.Sprintf("refresh-token"), refreshTokenAuth, &awsapigatewayv2integrations.HttpLambdaIntegrationProps{}),
	})
	api.AddRoutes(&awsapigatewayv2.AddRoutesOptions{
		Path:        jsii.Sprintf("/auth/revoke"),
		Methods:     &[]awsapigatewayv2.HttpMethod{awsapigatewayv2.HttpMethod_POST},
		Integration: awsapigatewayv2integrations.NewHttpLambdaIntegration(jsii.Sprintf("revoke-token"), revokeTokenAuth, &awsapigatewayv2integrations.HttpLambdaIntegrationProps{}),
	})

	props.CdkLambdaStackFunctions.downloadFileFunc.GrantInvoke(orchestrator.Role())
	props.CdkLambdaStackFunctions.uploadFileFunc.GrantInvoke(orchestrator.Role())
	props.CognitoStack.SetCallbackUrLs(jsii.Strings(callbackUri))
	props.TokenCodeTable.GrantReadWriteData(refreshTokenAuth)
	props.TokenCodeTable.GrantReadWriteData(tokenAuth)
	props.TokenCodeTable.GrantReadWriteData(callbackAuth)
	props.TokenCodeTable.GrantReadWriteData(revokeTokenAuth)

	awscdk.NewCfnOutput(stack, jsii.String("MoonenvApiGatewayUrl"), &awscdk.CfnOutputProps{Value: api.Url()})
}
