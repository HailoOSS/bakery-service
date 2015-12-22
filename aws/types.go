package aws

import (
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	log "github.com/cihub/seelog"
)

// Account structure
type Account struct {
	ID      string   `json:"id"`
	Regions []string `json:"regions"`
	SNSRole string   `json:"snsRole"`
}

// AssumeRole performs an API req to give temporary permissions to a service
func (a *Account) AssumeRole(sessionName string, duration time.Duration) (*aws.Config, error) {
	log.Debugf("Trying to assume role '%s' in '%s'", a.SNSRole, os.Getenv("EC2_REGION"))
	svc := sts.New(session.New(), &aws.Config{Region: aws.String(os.Getenv("EC2_REGION"))})

	return &aws.Config{
		Credentials: stscreds.NewCredentialsWithClient(svc, a.SNSRole),
		Region:      aws.String(os.Getenv("EC2_REGION")),
	}, nil
}
