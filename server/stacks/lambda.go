package stacks

import (
	"strconv"

	"github.com/aws/aws-cdk-go/awscdk"
	"github.com/aws/aws-cdk-go/awscdk/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/awslambdago"
	"github.com/aws/aws-cdk-go/awscdk/awss3"
	"github.com/aws/constructs-go/constructs/v3"
	"github.com/aws/jsii-runtime-go"
)

type CdkLambdaStackProps struct {
	awscdk.StackProps
	awss3.Bucket
	TokenCodeTable          awsdynamodb.Table
	TokenCodeStateIndexName *string
	AuthSubdomain           *string
	RestApiSubdomain        *string
}

type CdkLambdaStackFunctions struct {
	uploadFileFunc   awslambdago.GoFunction
	downloadFileFunc awslambdago.GoFunction
	tokenAuth        awslambdago.GoFunction
	callbackAuth     awslambdago.GoFunction
	refreshTokenAuth awslambdago.GoFunction
	revokeTokenAuth  awslambdago.GoFunction
	pullCommand      awslambdago.GoFunction
	pushCommand      awslambdago.GoFunction
}

func NewCdkLambdaStack(scope constructs.Construct, id string, props *CdkLambdaStackProps) *CdkLambdaStackFunctions {
	var sProps awscdk.StackProps

	if props != nil {
		sProps = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sProps)

	downloadFileFunc := awslambdago.NewGoFunction(stack, jsii.String("download-file-func"), &awslambdago.GoFunctionProps{
		MemorySize:   jsii.Number(128),
		Entry:        jsii.String("./lambdas/download-file"),
		Environment:  &map[string]*string{"S3Bucket": props.Bucket.BucketName()},
		FunctionName: jsii.String("moonenv-download-file"),
	})

	uploadFileFunc := awslambdago.NewGoFunction(stack, jsii.String("upload-file-func"), &awslambdago.GoFunctionProps{
		MemorySize:   jsii.Number(128),
		Entry:        jsii.String("./lambdas/upload-file"),
		Environment:  &map[string]*string{"S3Bucket": props.Bucket.BucketName()},
		FunctionName: jsii.String("moonenv-upload-file"),
	})

	tokenAuth := awslambdago.NewGoFunction(stack, jsii.String("MoonenvAuthToken"), &awslambdago.GoFunctionProps{
		MemorySize:   jsii.Number(128),
		Entry:        jsii.String("./lambdas/endpoints/auth/token"),
		FunctionName: jsii.String("moonenv-auth-token"),
		Environment: &map[string]*string{
			"TokenCodeTableName":       props.TokenCodeTable.TableName(),
			"PollingIntervalInSeconds": jsii.String(strconv.FormatInt(int64(3), 10)),
			"CognitoUrl":               props.AuthSubdomain,
			"CallbackUri":              GetApiGatewayCallbackUri(props.RestApiSubdomain),
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
			"CognitoUrl":         props.AuthSubdomain,
			"TokenCodeTableName": props.TokenCodeTable.TableName(),
		},
	})
	revokeTokenAuth := awslambdago.NewGoFunction(stack, jsii.Sprintf("MoonenvAuthRevokeToken"), &awslambdago.GoFunctionProps{
		MemorySize:   jsii.Number(128),
		Entry:        jsii.Sprintf("./lambdas/endpoints/auth/revoke"),
		FunctionName: jsii.Sprintf("moonenv-auth-revoke-token"),
		Environment: &map[string]*string{
			"CognitoUrl":         props.AuthSubdomain,
			"TokenCodeTableName": props.TokenCodeTable.TableName(),
		},
	})
	pullCommand := awslambdago.NewGoFunction(stack, jsii.String("MoonenvPullCommand"), &awslambdago.GoFunctionProps{
		MemorySize:   jsii.Number(128),
		Entry:        jsii.String("./lambdas/endpoints/orchestrator/pull"),
		FunctionName: jsii.String("moonenv-pull-command"),
		Environment: &map[string]*string{
			"AwsRegion":        props.StackProps.Env.Region,
			"DownloadFuncName": downloadFileFunc.FunctionArn(),
		},
	})
	pushCommand := awslambdago.NewGoFunction(stack, jsii.String("MoonenvPushCommand"), &awslambdago.GoFunctionProps{
		MemorySize:   jsii.Number(128),
		Entry:        jsii.String("./lambdas/endpoints/orchestrator/push"),
		FunctionName: jsii.String("moonenv-push-command"),
		Environment: &map[string]*string{
			"AwsRegion":      props.StackProps.Env.Region,
			"UploadFuncName": uploadFileFunc.FunctionArn(),
		},
	})

	props.Bucket.GrantRead(downloadFileFunc.Role(), nil)
	props.Bucket.GrantWrite(uploadFileFunc.Role(), nil)

	downloadFileFunc.GrantInvoke(pullCommand.Role())
	uploadFileFunc.GrantInvoke(pushCommand.Role())

	authTypes := []awslambdago.GoFunction{refreshTokenAuth, tokenAuth, callbackAuth, revokeTokenAuth}

	for _, auth := range authTypes {
		props.TokenCodeTable.GrantReadWriteData(auth)
	}

	return &CdkLambdaStackFunctions{
		uploadFileFunc:   uploadFileFunc,
		downloadFileFunc: downloadFileFunc,
		tokenAuth:        tokenAuth,
		callbackAuth:     callbackAuth,
		refreshTokenAuth: refreshTokenAuth,
		revokeTokenAuth:  revokeTokenAuth,
		pullCommand:      pullCommand,
		pushCommand:      pushCommand,
	}
}
