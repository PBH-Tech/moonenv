package stacks

import (
	"strconv"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3assets"
	"github.com/aws/constructs-go/constructs/v10"
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
	uploadFileFunc   awslambda.Function
	downloadFileFunc awslambda.Function
	tokenAuth        awslambda.Function
	callbackAuth     awslambda.Function
	refreshTokenAuth awslambda.Function
	revokeTokenAuth  awslambda.Function
	pullCommand      awslambda.Function
	pushCommand      awslambda.Function
}

func NewCdkLambdaStack(scope constructs.Construct, id string, props *CdkLambdaStackProps) *CdkLambdaStackFunctions {
	var sProps awscdk.StackProps

	if props != nil {
		sProps = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sProps)

	downloadFileFunc := awslambda.NewFunction(stack, jsii.String("download-file-func"), &awslambda.FunctionProps{
		Runtime:      awslambda.Runtime_GO_1_X(),
		MemorySize:   jsii.Number(128),
		Code:         awslambda.AssetCode_FromAsset(jsii.String("./lambdas/download-file"), &awss3assets.AssetOptions{}),
		Environment:  &map[string]*string{"S3Bucket": props.Bucket.BucketName()},
		FunctionName: jsii.String("moonenv-download-file"),
		Handler:      jsii.String("main.handler"),
	})

	uploadFileFunc := awslambda.NewFunction(stack, jsii.String("upload-file-func"), &awslambda.FunctionProps{
		Runtime:      awslambda.Runtime_GO_1_X(),
		MemorySize:   jsii.Number(128),
		Code:         awslambda.AssetCode_FromAsset(jsii.String("./lambdas/upload-file"), &awss3assets.AssetOptions{}),
		Environment:  &map[string]*string{"S3Bucket": props.Bucket.BucketName()},
		FunctionName: jsii.String("moonenv-upload-file"),
		Handler:      jsii.String("main.handler"),
	})

	tokenAuth := awslambda.NewFunction(stack, jsii.String("MoonenvAuthToken"), &awslambda.FunctionProps{
		Runtime:      awslambda.Runtime_GO_1_X(),
		MemorySize:   jsii.Number(128),
		Code:         awslambda.AssetCode_FromAsset(jsii.String("./lambdas/endpoints/auth/token"), &awss3assets.AssetOptions{}),
		FunctionName: jsii.String("moonenv-auth-token"),
		Environment: &map[string]*string{
			"TokenCodeTableName":       props.TokenCodeTable.TableName(),
			"PollingIntervalInSeconds": jsii.String(strconv.FormatInt(int64(3), 10)),
			"CognitoUrl":               props.AuthSubdomain,
			"CallbackUri":              GetApiGatewayCallbackUri(props.RestApiSubdomain),
		},
		Handler: jsii.String("main.handler"),
	})
	callbackAuth := awslambda.NewFunction(stack, jsii.Sprintf("MoonenvAuthCallback"), &awslambda.FunctionProps{
		Runtime:      awslambda.Runtime_GO_1_X(),
		MemorySize:   jsii.Number(128),
		Code:         awslambda.AssetCode_FromAsset(jsii.Sprintf("./lambdas/endpoints/auth/callback"), &awss3assets.AssetOptions{}),
		FunctionName: jsii.Sprintf("moonenv-auth-callback"),
		Environment: &map[string]*string{
			"StateIndexName":     props.TokenCodeStateIndexName,
			"TokenCodeTableName": props.TokenCodeTable.TableName(),
		},
		Handler: jsii.String("main.handler"),
	})
	refreshTokenAuth := awslambda.NewFunction(stack, jsii.Sprintf("MoonenvAuthRefreshToken"), &awslambda.FunctionProps{
		Runtime:      awslambda.Runtime_GO_1_X(),
		MemorySize:   jsii.Number(128),
		Code:         awslambda.AssetCode_FromAsset(jsii.Sprintf("./lambdas/endpoints/auth/refresh"), &awss3assets.AssetOptions{}),
		FunctionName: jsii.Sprintf("moonenv-auth-refresh-token"),
		Environment: &map[string]*string{
			"CognitoUrl":         props.AuthSubdomain,
			"TokenCodeTableName": props.TokenCodeTable.TableName(),
		},
		Handler: jsii.String("main.handler"),
	})
	revokeTokenAuth := awslambda.NewFunction(stack, jsii.Sprintf("MoonenvAuthRevokeToken"), &awslambda.FunctionProps{
		Runtime:      awslambda.Runtime_GO_1_X(),
		MemorySize:   jsii.Number(128),
		Code:         awslambda.AssetCode_FromAsset(jsii.Sprintf("./lambdas/endpoints/auth/revoke"), &awss3assets.AssetOptions{}),
		FunctionName: jsii.Sprintf("moonenv-auth-revoke-token"),
		Environment: &map[string]*string{
			"CognitoUrl":         props.AuthSubdomain,
			"TokenCodeTableName": props.TokenCodeTable.TableName(),
		},
		Handler: jsii.String("main.handler"),
	})
	pullCommand := awslambda.NewFunction(stack, jsii.String("MoonenvPullCommand"), &awslambda.FunctionProps{
		Runtime:      awslambda.Runtime_GO_1_X(),
		MemorySize:   jsii.Number(128),
		Code:         awslambda.AssetCode_FromAsset(jsii.String("./lambdas/endpoints/orchestrator/pull"), &awss3assets.AssetOptions{}),
		FunctionName: jsii.String("moonenv-pull-command"),
		Environment: &map[string]*string{
			"AwsRegion":        props.StackProps.Env.Region,
			"DownloadFuncName": downloadFileFunc.FunctionArn(),
		},
		Handler: jsii.String("main.handler"),
	})
	pushCommand := awslambda.NewFunction(stack, jsii.String("MoonenvPushCommand"), &awslambda.FunctionProps{
		Runtime:      awslambda.Runtime_GO_1_X(),
		MemorySize:   jsii.Number(128),
		Code:         awslambda.AssetCode_FromAsset(jsii.String("./lambdas/endpoints/orchestrator/push"), &awss3assets.AssetOptions{}),
		FunctionName: jsii.String("moonenv-push-command"),
		Environment: &map[string]*string{
			"AwsRegion":      props.StackProps.Env.Region,
			"UploadFuncName": uploadFileFunc.FunctionArn(),
		},
		Handler: jsii.String("main.handler"),
	})

	props.Bucket.GrantRead(downloadFileFunc.Role(), nil)
	props.Bucket.GrantWrite(uploadFileFunc.Role(), map[string]string{}, nil)

	downloadFileFunc.GrantInvoke(pullCommand.Role())
	uploadFileFunc.GrantInvoke(pushCommand.Role())

	authTypes := []awslambda.Function{refreshTokenAuth, tokenAuth, callbackAuth, revokeTokenAuth}

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
