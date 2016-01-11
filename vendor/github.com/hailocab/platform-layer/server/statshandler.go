package server

import (
	"github.com/hailocab/protobuf/proto"

	"github.com/hailocab/platform-layer/errors"
	"github.com/hailocab/platform-layer/stats"
)

// registerStats starts the runtime stats collection and monitoring
func registerStats() {
	stats.ServiceName = Name
	stats.ServiceVersion = Version
	stats.ServiceType = "h2.platform"
	stats.InstanceID = InstanceID
	for _, ep := range reg.iterate() {
		stats.RegisterEndpoint(ep)
	}

	stats.Start()
}

// statsHandler handles inbound requests to `stats` endpoint
func statsHandler(req *Request) (proto.Message, errors.Error) {
	return stats.Get(), nil
}
