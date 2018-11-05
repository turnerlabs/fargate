package sts

import (
	"github.com/aws/aws-sdk-go/aws/session"
	awssts "github.com/aws/aws-sdk-go/service/sts"
	"github.com/turnerlabs/fargate/console"
)

//STS represents an STS API
type STS struct {
	svc *awssts.STS
}

//New creates a new STS value
func New(sess *session.Session) STS {
	return STS{
		svc: awssts.New(sess),
	}
}

// CallerIdentity Contains the response to a successful GetCallerIdentity request, including
// information about the entity making the request.
type CallerIdentity struct {
	Account string
	ARN     string
	UserID  string
}

//GetCallerIdentity calls GetCallerIdentity
func (s *STS) GetCallerIdentity() CallerIdentity {
	input := &awssts.GetCallerIdentityInput{}
	resp, err := s.svc.GetCallerIdentity(input)
	if err != nil {
		console.ErrorExit(err, "Error calling GetCallerIdentity")
	}
	result := CallerIdentity{
		Account: *resp.Account,
		ARN:     *resp.Arn,
		UserID:  *resp.UserId,
	}

	return result
}
