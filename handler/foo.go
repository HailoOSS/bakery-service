package handler

import (
	"code.google.com/p/goprotobuf/proto"
	log "github.com/cihub/seelog"
	"github.com/hailocab/go-platform-layer/errors"
	"github.com/hailocab/go-platform-layer/server"
	foo "github.com/hailocab/{{REPONAME}}/proto/foo"
)

// Foo does <what>? Remember to add Godocs to your handlers, and follow the Go convention of starting with the function name
func Foo(req *server.Request) (proto.Message, errors.Error) {
	log.Infof("Doing foo %+v", req)

	request := &foo.Request{}
	if err := req.Unmarshal(request); err != nil {
		return nil, errors.BadRequest(server.Name+".foo", err.Error())
	}

	// we probably want to make use of the request parameter that we know we will be passed:
	log.Debugf("Received bar=%v", request.GetBar())

	// INSERT CODE HERE TO ACTUALLY DO SOMETHING!

	// then we can make a response
	rsp := &foo.Response{
		Baz: proto.String("This is what we return"),
	}

	return rsp, nil
}
