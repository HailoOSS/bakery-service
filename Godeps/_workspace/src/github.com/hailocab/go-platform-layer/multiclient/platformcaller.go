package multiclient

import (
	"github.com/hailocab/protobuf/proto"

	"github.com/hailocab/go-platform-layer/client"
	"github.com/hailocab/go-platform-layer/errors"
)

// PlatformCaller is the default caller and makes requests via the platform layer
// RPC mechanism (eg: RabbitMQ)
func PlatformCaller() Caller {
	return func(req *client.Request, rsp proto.Message) errors.Error {
		return client.Req(req, rsp)
	}
}
