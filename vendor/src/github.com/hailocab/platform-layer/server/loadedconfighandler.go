package server

import (
	"github.com/hailocab/protobuf/proto"

	"github.com/hailocab/platform-layer/errors"
	"github.com/hailocab/service-layer/config"

	loadedconfigproto "github.com/hailocab/platform-layer/proto/loadedconfig"
)

// loadedConfigHandler handles inbound requests to `loadedconfig` endpoint
func loadedConfigHandler(req *Request) (proto.Message, errors.Error) {
	configJson := string(config.Raw())
	return &loadedconfigproto.Response{
		Config: proto.String(configJson),
	}, nil
}
