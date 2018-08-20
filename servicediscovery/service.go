package servicediscovery

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	awssd "github.com/aws/aws-sdk-go/service/servicediscovery"
	"github.com/turnerlabs/fargate/console"
)

type DnsRecord struct {
	TTL  int64
	Type string
}

type Service struct {
	Id         string
	DnsRecords []DnsRecord
	Name       string
	Namespace  Namespace
}

func (sd *ServiceDiscovery) GetService(registryArn string) Service {
	var service Service

	arnSlice := strings.Split(registryArn, "/")
	id := arnSlice[len(arnSlice)-1]

	resp, err := sd.svc.GetService(
		&awssd.GetServiceInput{
			Id: aws.String(id),
		},
	)

	if err != nil {
		console.ErrorExit(err, "Could not describe ServiceDiscovery service")
	}

	namespace := sd.GetNamespace(
		aws.StringValue(resp.Service.DnsConfig.NamespaceId),
	)

	service = Service{
		Id:        id,
		Name:      aws.StringValue(resp.Service.Name),
		Namespace: namespace,
	}

	for _, dnsRecord := range resp.Service.DnsConfig.DnsRecords {
		service.DnsRecords = append(
			service.DnsRecords,
			DnsRecord{
				TTL:  aws.Int64Value(dnsRecord.TTL),
				Type: aws.StringValue(dnsRecord.Type),
			},
		)
	}

	return service
}
