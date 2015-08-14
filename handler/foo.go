package handler

import (
	log "github.com/cihub/seelog"
	"github.com/hailocab/protobuf/proto"

	foo "github.com/hailocab/bakery-service/proto/foo"
	"github.com/hailocab/go-platform-layer/errors"
	"github.com/hailocab/go-platform-layer/server"
)

// Foo does <what>? Remember to add Godocs to your handlers, and follow the Go convention of starting with the function name
func Foo(req *server.Request) (proto.Message, errors.Error) {
	log.Infof("Doing foo %+v", req)

	request := req.Data().(*foo.Request)

	// we probably want to make use of the request parameter that we know we will be passed:
	log.Debugf("Received bar=%v", request.GetBar())

	// INSERT CODE HERE TO ACTUALLY DO SOMETHING!

	// then we can make a response
	rsp := &foo.Response{
		Baz: proto.String("This is what we return"),
	}

	return rsp, nil
}
