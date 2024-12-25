package stacks

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscertificatemanager"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsroute53"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type CdkRoute53StackProps struct {
	awscdk.StackProps
	HostZoneId              *string
	MoonenvDomain           *string
	CertificateArnInUsEast1 *string
}

type CdkRoute53StackResource struct {
	awsroute53.IHostedZone
	awscertificatemanager.Certificate
	UsEast1Certificate awscertificatemanager.ICertificate
}

func NewRoute53Stack(scope constructs.Construct, id string, props *CdkRoute53StackProps) CdkRoute53StackResource {
	var sProps awscdk.StackProps

	if props != nil {
		sProps = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sProps)

	hostZone := awsroute53.HostedZone_FromHostedZoneAttributes(stack, jsii.String("MoonenvHostZone"), &awsroute53.HostedZoneAttributes{
		HostedZoneId: props.HostZoneId,
		ZoneName:     props.MoonenvDomain,
	})
	certificate := awscertificatemanager.NewCertificate(stack, jsii.String("MoonenvCertificate"), &awscertificatemanager.CertificateProps{
		DomainName: jsii.Sprintf("*.%s", *props.MoonenvDomain),
		Validation: awscertificatemanager.CertificateValidation_FromDns(hostZone),
	})
	usEast1Certificate := awscertificatemanager.Certificate_FromCertificateArn(stack, jsii.String("MoonenvUsEat1Certificate"), props.CertificateArnInUsEast1)

	return CdkRoute53StackResource{
		IHostedZone:        hostZone,
		Certificate:        certificate,
		UsEast1Certificate: usEast1Certificate,
	}

}
