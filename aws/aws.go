package aws

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"

	"github.com/hailocab/go-service-layer/config"

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

// UploadPart takes a slice of data and appends it to a file
func UploadPart(bucket string, path string, part int64, body io.ReadSeeker) error {
	config, err := Auth(DefaultAccount)
	if err != nil {
		return fmt.Errorf("Unable to auth: %v", err)
	}

	log.Debugf("Trying to save part %d to: '%s/%s'", part, bucket, path)
	svc := s3.New(session.New(), config)

	resp, err := svc.UploadPart(&s3.UploadPartInput{
		Bucket:     aws.String(bucket),
		Key:        aws.String(path),
		PartNumber: aws.Int64(part),
		Body:       body,
		UploadId:   aws.String("bakery"),
	})

	if err != nil {
		return fmt.Errorf("Unable to save part %d for %q", part, path)
	}

	log.Debugf("UploadPart resp: %v", resp)

	return nil
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

func LoadEncryptedAccountInfo() (map[string]string, error) {
	credentials, err := config.AtPath(
		"hailo",
		"service",
		"bakery",
		"credentials",
	).Decrypt()

	if err != nil {
		return map[string]string{}, err
	}

	return map[string]string{
		"aws_access_key_id":     credentials.AtPath("aws_access_key_id").AsString(),
		"aws_secret_access_key": credentials.AtPath("aws_secret_access_key").AsString(),
	}, nil
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
