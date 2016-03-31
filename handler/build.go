package handler

import (
	"fmt"
	"os"
	"path/filepath"

	protoBuild "github.com/hailocab/bakery-service/proto/build"

	"github.com/hailocab/bakery-service/aws"
	"github.com/hailocab/bakery-service/elastic"
	"github.com/hailocab/bakery-service/packer"
	"github.com/hailocab/bakery-service/packer/ui"

	"github.com/hailocab/platform-layer/errors"
	"github.com/hailocab/platform-layer/server"

	log "github.com/cihub/seelog"
	"github.com/hailocab/protobuf/proto"
	"github.com/hashicorp/go-getter"
	"github.com/nu7hatch/gouuid"
)

const (
	// BuildEndpoint name of endpoint
	BuildEndpoint = "com.hailocab.infrastructure.bakery.build"

	// BucketName S3 bucket to find templates
	BucketName = "hailo-bakery"

	// BucketLogPath path to store logs
	BucketLogPath = "logs"

	// BucketTemplatePath path where templates are stored
	BucketTemplatePath = "templates"
)

// Build endpoint
func Build(req *server.Request) (proto.Message, errors.Error) {
	var (
		p   *packer.Packer
		err error
	)

	request := req.Data().(*protoBuild.Request)
	reqVars := map[string]string{}
	for _, v := range request.GetVariables() {
		reqVars[v.GetKey()] = v.GetValue()
	}

	template := request.GetTemplate()
	log.Infof("Requested Template: %v", template)

	dir, err := packer.TemporaryDir()
	if err != nil {
		return nil, errors.InternalServerError(BuildEndpoint, err)
	}

	s3target := fmt.Sprintf("s3://%s/%s/%s.zip", BucketName, BucketTemplatePath, template)
	log.Infof("s3: %s", s3target)

	if err := getter.GetFile(dir, s3target); err != nil {
		return nil, errors.InternalServerError(BuildEndpoint, err)
	}

	rc, err := os.Open(filepath.Join(dir, fmt.Sprintf("%s.json", template)))

	id, err := uuid.NewV4()
	if err != nil {
		return nil, errors.InternalServerError(BuildEndpoint,
			fmt.Sprintf("Unable to create ID: %v", err),
		)
	}

	e, err := elastic.NewWithDefaults()
	if err != nil {
		return nil, errors.InternalServerError(BuildEndpoint,
			fmt.Sprintf("Unable to create new elastic: %v", err),
		)
	}

	ui := ui.New(
		ui.AddCaller("echo", &ui.EchoCaller{}),
		ui.AddCaller("elastic", ui.NewElasticCaller(id.String(), e)),
	)

	p, err = packer.New(rc, ui)
	if err != nil {
		return nil, errors.InternalServerError(BuildEndpoint,
			fmt.Sprintf("Can't build resource: %v", err),
		)
	}

	config, err := aws.Auth(aws.DefaultAccount)
	if err != nil {
		return nil, errors.InternalServerError(BuildEndpoint,
			fmt.Sprintf("Unable to get AWS configuration"),
		)
	}

	creds, err := config.Credentials.Get()
	if err != nil {
		return nil, errors.InternalServerError(BuildEndpoint,
			fmt.Sprintf("Unable to get AWS credentials"),
		)
	}

	vars := packer.ExtractVariables(p.Template.Variables, map[string]string{
		"cwd":                   dir,
		"aws_access_key_id":     creds.AccessKeyID,
		"aws_secret_access_key": creds.SecretAccessKey,
		"aws_session_token":     creds.SessionToken,
	})

	go p.Build(vars)

	return &protoBuild.Response{
		Id: proto.String(id.String()),
	}, nil
}
