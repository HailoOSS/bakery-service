package aws

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"

	"github.com/hailocab/service-layer/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	log "github.com/cihub/seelog"
)

var (
	// DefaultAccount ID of default account to operate on
	DefaultAccount = "864806739507"

	accounts []Account
)

// Init triggers a config load
func Init() {
	var err error

	accounts, err = loadAccountInfo()
	if err != nil {
		panic(err)
	}
}

// Auth assumes a role specified
func Auth(accountID string) (*aws.Config, error) {
	var account Account
	for _, a := range accounts {
		if a.ID == accountID {
			account = a
			break
		}
	}

	return account.AssumeRole(randString(), 3600)
}

// Credentials returns a credentials struct
func Credentials() (credentials.Value, error) {
	config, err := Auth(DefaultAccount)
	if err != nil {
		return credentials.Value{}, fmt.Errorf("Unable to auth: %v", err)
	}

	value, err := config.Credentials.Get()
	if err != nil {
		return credentials.Value{}, fmt.Errorf("Unable to get credentials: %v", err)
	}

	return value, nil
}

// GetS3Object returns an object from s3
func GetS3Object(bucket string, key string) (io.ReadCloser, error) {
	config, err := Auth(DefaultAccount)
	if err != nil {
		return nil, fmt.Errorf("Unable to auth: %v", err)
	}

	svc := s3.New(session.New(), config)

	resp, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		return nil, fmt.Errorf("Unable to fetch '%s/%s': %v", bucket, key, err)
	}

	return resp.Body, nil
}

func loadAccountInfo() ([]Account, error) {
	accountConfig := config.AtPath(
		"hailo",
		"service",
		"bakery",
		"accounts",
	).AsJson()

	log.Debugf("Bakery accounts: %v", accountConfig)

	var accs []Account
	if err := json.Unmarshal(accountConfig, &accs); err != nil {
		return nil, err
	}

	log.Debugf("Actual accounts: %v", accs)
	return accs, nil
}

func randString() string {
	alphanum := "0123456789abcdefghigklmnopqrst"

	var bytes = make([]byte, 10)
	rand.Read(bytes)

	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}

	return string(bytes)
}