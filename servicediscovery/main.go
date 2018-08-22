package servicediscovery

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/servicediscovery"
)

type ServiceDiscovery struct {
	svc *servicediscovery.ServiceDiscovery
}

func New(sess *session.Session) ServiceDiscovery {
	return ServiceDiscovery{
		svc: servicediscovery.New(sess),
	}
}
