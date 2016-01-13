package handler

import (
	"fmt"

	protoBuild "github.com/hailocab/bakery-service/proto/build"

	"github.com/hailocab/bakery-service/aws"
	"github.com/hailocab/bakery-service/packer"

	"github.com/hailocab/platform-layer/errors"
	"github.com/hailocab/platform-layer/server"

	log "github.com/cihub/seelog"
	"github.com/hailocab/protobuf/proto"
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

	rc, err := aws.GetS3Object(BucketName, fmt.Sprintf("%s/%s.json", BucketTemplatePath, template))
	if err != nil {
		return nil, errors.BadRequest(BuildEndpoint,
			fmt.Sprintf("Unable to get object: %v", err),
		)
	}

	p, err = packer.New(rc, BucketName, BucketLogPath)
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
		"aws_access_key_id":     creds.AccessKeyID,
		"aws_secret_access_key": creds.SecretAccessKey,
		"aws_session_token":     creds.SessionToken,
	})

	go p.Build(vars)

	return &protoBuild.Response{
		Id: proto.String("lolz"),
	}, nil
}
