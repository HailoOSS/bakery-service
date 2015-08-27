package handler

import (
	"fmt"

	protoBuild "github.com/hailocab/bakery-service/proto/build"

	"github.com/hailocab/bakery-service/aws"
	"github.com/hailocab/bakery-service/packer"

	"github.com/hailocab/go-platform-layer/errors"
	"github.com/hailocab/go-platform-layer/server"

	log "github.com/cihub/seelog"
	"github.com/hailocab/protobuf/proto"
)

const (
	BuildEndpoint  = "com.hailocab.infrastructure.bakery.build"
	TemplateBucket = "hailo-bakery"
)

func Build(req *server.Request) (proto.Message, errors.Error) {
	request := req.Data().(*protoBuild.Request)

	template := request.GetTemplate()
	log.Infof("Requested Template: %v", template)

	rc, err := aws.GetS3Object(TemplateBucket, template)
	if err != nil {
		return nil, errors.BadRequest(BuildEndpoint,
			fmt.Sprintf("Unable to get object: %v", err),
		)
	}

	packer.Build(rc)

	return &protoBuild.Response{
		Id: proto.String("lolz"),
	}, nil
}
