package stacks

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type CdkTableStackProps struct {
	awscdk.StackProps
	TableName    string
	PartitionKey awsdynamodb.Attribute
}

func NewTableStack(scope constructs.Construct, id string, props *CdkTableStackProps) awsdynamodb.Table {
	var sProps awscdk.StackProps

	if props != nil {
		sProps = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sProps)

	return awsdynamodb.NewTable(stack, jsii.String("MoonenvTokenCode"), &awsdynamodb.TableProps{
		TableName:    &props.TableName,
		PartitionKey: &props.PartitionKey,
	})
}
