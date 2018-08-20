package servicediscovery

import (
	"github.com/aws/aws-sdk-go/aws"
	awssd "github.com/aws/aws-sdk-go/service/servicediscovery"
	"github.com/turnerlabs/fargate/console"
)

type Namespace struct {
	Id      string
	Name    string
	Private bool
}

func (sd *ServiceDiscovery) GetNamespace(namespaceId string) Namespace {
	var namespace Namespace

	resp, err := sd.svc.GetNamespace(
		&awssd.GetNamespaceInput{
			Id: aws.String(namespaceId),
		},
	)

	if err != nil {
		console.ErrorExit(err, "Could not describe ServiceDiscovery namespace")
	}

	namespace = Namespace{
		Id:      namespaceId,
		Name:    aws.StringValue(resp.Namespace.Name),
		Private: aws.StringValue(resp.Namespace.Type) == awssd.NamespaceTypeDnsPrivate,
	}

	return namespace
}
