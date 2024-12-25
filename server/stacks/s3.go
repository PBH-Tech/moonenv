package stacks

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type CdkS3StackProps struct {
	awscdk.StackProps
}

func NewS3BucketStack(scope constructs.Construct, id string, props *CdkS3StackProps) awss3.Bucket {
	var sProps awscdk.StackProps

	if props != nil {
		sProps = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sProps)

	bucket := awss3.NewBucket(stack, jsii.String("moonenv-bucket"), &awss3.BucketProps{
		BucketName: jsii.String("moonenv-bucket"),
		Versioned:  jsii.Bool(true),
	})

	return bucket
}
