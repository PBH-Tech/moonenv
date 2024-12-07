package stacks

import (
	"errors"

	"github.com/aws/aws-cdk-go/awscdk"
	"github.com/aws/aws-cdk-go/awscdk/awslambdago"
	"github.com/aws/aws-cdk-go/awscdk/awss3"
	"github.com/aws/constructs-go/constructs/v3"
	"github.com/aws/jsii-runtime-go"
)

type CdkLambdaStackProps struct {
	awscdk.StackProps
	awss3.Bucket
}

type CdkLambdaStackFunctions struct {
	uploadFileFunc   awslambdago.GoFunction
	downloadFileFunc awslambdago.GoFunction
}

func NewCdkLambdaStack(scope constructs.Construct, id string, props *CdkLambdaStackProps) (*CdkLambdaStackFunctions, error) {
	var sProps awscdk.StackProps

	if props != nil {
		sProps = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sProps)

	if props.Bucket == nil {
		return nil, errors.New("BUCKET SHOULD BE DEFINED")
	}

	downloadFileFunc := awslambdago.NewGoFunction(stack, jsii.String("download-file-func"), &awslambdago.GoFunctionProps{
		MemorySize:  jsii.Number(128),
		Entry:       jsii.String("./lambdas/download-file"),
		Environment: &map[string]*string{"S3Bucket": props.Bucket.BucketName()},
	})

	uploadFileFunc := awslambdago.NewGoFunction(stack, jsii.String("upload-file-func"), &awslambdago.GoFunctionProps{
		MemorySize:  jsii.Number(128),
		Entry:       jsii.String("./lambdas/upload-file"),
		Environment: &map[string]*string{"S3Bucket": props.Bucket.BucketName()},
	})

	// Grant read permissions to the download functioni
	props.Bucket.GrantRead(downloadFileFunc.Role(), nil)

	// Grant write permissions to the upload function
	props.Bucket.GrantWrite(uploadFileFunc.Role(), nil)

	return &CdkLambdaStackFunctions{uploadFileFunc: uploadFileFunc, downloadFileFunc: downloadFileFunc}, nil
}
