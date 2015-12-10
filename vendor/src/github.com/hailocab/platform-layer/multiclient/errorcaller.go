package multiclient

import (
	"github.com/hailocab/protobuf/proto"

	"github.com/hailocab/platform-layer/client"
	"github.com/hailocab/platform-layer/errors"
)

// ErrorCaller is a very simple caller that just returns an error
// If no error provided, defaults to a `NotFound` error with code "errorcaller.notfound"
func ErrorCaller(err errors.Error) Caller {
	return func(req *client.Request, rsp proto.Message) errors.Error {
		if err != nil {
			return err
		}
		return errors.NotFound("errorcaller.notfound", "No error supplied.")
	}
}
