package aws

import (
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/service/sts"
	log "github.com/cihub/seelog"
)

type Account struct {
	Id      string   `json:"id"`
	Regions []string `json:"regions"`
	SNSRole string   `json:"snsRole"`
}

func (a *Account) AssumeRole(session string, duration time.Duration) (*aws.Config, error) {
	log.Debugf("Trying to assume role '%s' in '%s'", a.SNSRole, os.Getenv("EC2_REGION"))
	svc := sts.New(&aws.Config{Region: aws.String(os.Getenv("EC2_REGION"))})

	return &aws.Config{
		Credentials: stscreds.NewCredentials(svc, a.SNSRole, duration),
		Region:      aws.String(os.Getenv("EC2_REGION")),
	}, nil
}
