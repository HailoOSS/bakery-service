package server

import (
	"github.com/hailocab/protobuf/proto"

	"github.com/hailocab/platform-layer/errors"
	"github.com/hailocab/platform-layer/healthcheck"
)

// healthHandler handles inbound requests to `health` endpoint
func healthHandler(req *Request) (proto.Message, errors.Error) {
	return healthcheck.Status(), nil
}
