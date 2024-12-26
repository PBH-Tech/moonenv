package schema

import (
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/jsii-runtime-go"
)

var (
	PushCommandRequestSchema = awsapigateway.JsonSchema{
		Type:     awsapigateway.JsonSchemaType_OBJECT,
		Required: &[]*string{jsii.String("b64String")},
		Properties: &map[string]*awsapigateway.JsonSchema{
			"b64String": {
				Type: awsapigateway.JsonSchemaType_STRING,
			},
		},
	}
)
