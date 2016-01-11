package server

import (
	"github.com/hailocab/protobuf/proto"

	"github.com/hailocab/platform-layer/client"
)

var defaultScoper = Scoper()

// Scoper mints something that is able to yield a scoped request for this server
func Scoper() *serverScoper {
	return &serverScoper{}
}

type serverScoper struct{}

// ScopedRequest yields a new request from our server's scope (setting "from" details)
func (ss *serverScoper) ScopedRequest(service, endpoint string, payload proto.Message) (*client.Request, error) {
	return ScopedRequest(service, endpoint, payload)
}

// Context returns some context to base error messages from, eg: server name
func (ss *serverScoper) Context() string {
	return Name
}
