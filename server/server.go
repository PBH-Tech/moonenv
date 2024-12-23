package main

import (
	"errors"
	"os"

	"github.com/PBH-Tech/moonenv/stacks"
	"github.com/aws/aws-cdk-go/awscdk"
	"github.com/aws/aws-cdk-go/awscdk/awsdynamodb"
	"github.com/aws/jsii-runtime-go"
)

type ServerStackProps struct {
	awscdk.StackProps
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	bucket := stacks.NewS3BucketStack(app, "MoonenvS3Stack", &stacks.CdkS3StackProps{
		StackProps: awscdk.StackProps{
			Env:       env(),
			StackName: jsii.String("moonenv-s3"),
		}})
	lambdas, err := stacks.NewCdkLambdaStack(app, "MoonenvLambdaStack", &stacks.CdkLambdaStackProps{
		StackProps: awscdk.StackProps{
			Env:       env(),
			StackName: jsii.String("moonenv-lambda"),
		},
		Bucket: bucket,
	})

	if err != nil {
		errors.New(err.Error())
	}

	tokenCodeTable := stacks.NewTableStack(app, "MoonenvTokenCodeTable", &stacks.CdkTableStackProps{
		StackProps: awscdk.StackProps{
			Env:       env(),
			StackName: jsii.String("moonenv-token-code-table"),
		},
		TableName:    *jsii.String("moonenv-token-code"),
		PartitionKey: awsdynamodb.Attribute{Name: jsii.String("deviceCode"), Type: awsdynamodb.AttributeType_STRING},
	})

	tokenCodeStateIndexName := jsii.Sprintf("state-index")
	tokenCodeTable.AddGlobalSecondaryIndex(&awsdynamodb.GlobalSecondaryIndexProps{
		IndexName: tokenCodeStateIndexName,
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String("state"),
			Type: awsdynamodb.AttributeType_STRING,
		},
	})

	cognitoStack := stacks.NewCognitoStack(app, "MoonenvCognitoStack", &stacks.CdkCognitoStackProps{
		StackProps: awscdk.StackProps{
			Env:       env(),
			StackName: jsii.String("moonenv-cognito"),
		},
	})

	stacks.NewApiGatewayStack(app, "MoonenvApiGatewayStack", &stacks.CdkApiGatewayProps{
		StackProps: awscdk.StackProps{
			Env:       env(),
			StackName: jsii.String("moonenv-api-gateway"),
		},
		CdkLambdaStackFunctions: *lambdas,
		TokenCodeTable:          tokenCodeTable,
		CognitoStack:            *cognitoStack,
		TokenCodeStateIndexName: tokenCodeStateIndexName,
	})

	app.Synth(nil)

}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	// If unspecified, this stack will be "environment-agnostic".
	// Account/Region-dependent features and context lookups will not work, but a
	// single synthesized template can be deployed anywhere.
	//---------------------------------------------------------------------------
	// return nil

	// Uncomment if you know exactly what account and region you want to deploy
	// the stack to. This is the recommendation for production stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String("123456789012"),
	//  Region:  jsii.String("us-east-1"),
	// }

	// Uncomment to specialize this stack for the AWS Account and Region that are
	// implied by the current CLI configuration. This is recommended for dev
	// stacks.
	//---------------------------------------------------------------------------
	return &awscdk.Environment{
		Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
		Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	}
}
