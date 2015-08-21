package server

import (
	"fmt"
	"io"
	"time"

	log "github.com/cihub/seelog"
	"github.com/hailocab/protobuf/proto"

	"github.com/hailocab/go-platform-layer/errors"
)

// Handler interface
type Handler func(req *Request) (proto.Message, errors.Error)

// PostConnectHandler repesents a function we fun after connecting to RabbitMQ
type PostConnectHandler func()

// RegisterPostConnectHandler adds a post connect handler to the map so we can run it
func RegisterPostConnectHandler(pch PostConnectHandler) {
	// Don't see the need in locking this. Feel free to add if you think it's needed
	postConnHdlrs = append(postConnHdlrs, pch)
}

// commonLogHandler will log to w using the Apache common log format
// http://httpd.apache.org/docs/2.2/logs.html#common
// If w is nil, nothing will be logged
func commonLogHandler(w io.Writer, h Handler) Handler {
	if w == nil {
		return h
	}

	return func(req *Request) (proto.Message, errors.Error) {
		var userId string
		if req.Auth() != nil && req.Auth().AuthUser() != nil {
			userId = req.Auth().AuthUser().Id
		}

		var err errors.Error
		var m proto.Message

		// In defer in case the handler panics
		defer func() {
			status := uint32(200)
			if err != nil {
				status = err.HttpCode()
			}

			size := 0
			if m != nil {
				log.Debug(m.String())
				size = len(m.String())
			}

			fmt.Fprintf(w, "%s - %s [%s] \"%s %s %s\" %d %d\n",
				req.From(),
				userId,
				time.Now().Format("02/Jan/2006:15:04:05 -0700"),
				"GET", // Treat them all as GET's at the moment
				req.Endpoint(),
				"HTTP/1.0", // Has to be HTTP or apachetop ignores it
				status,
				size,
			)
		}()

		// Execute the actual handler
		m, err = h(req)
		return m, err
	}
}
