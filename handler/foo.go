package handler

import (
	"code.google.com/p/goprotobuf/proto"
	"fmt"
	log "github.com/cihub/seelog"
	"github.com/hailocab/go-platform-layer/errors"
	"github.com/hailocab/go-platform-layer/server"
	foo "github.com/hailocab/{{REPONAME}}/proto/foo"
)

func Foo(req *server.Request) (proto.Message, *errors.Error) {
	log.Infof("Doing foo %+v", req)

	request := &foo.Request{}
	if err := req.Unmarshal(request); err != nil {
		return nil, errors.InternalServerError("com.hailocab.service.{{SERVICENAME}}.foo", fmt.Sprintf("%v", err.Error()))
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
