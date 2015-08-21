package server

import (
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/hailocab/protobuf/proto"

	perrors "github.com/hailocab/go-platform-layer/errors"
	traceproto "github.com/hailocab/go-platform-layer/proto/trace"
	"github.com/hailocab/go-platform-layer/stats"
	inst "github.com/hailocab/go-service-layer/instrumentation"
	trace "github.com/hailocab/go-service-layer/trace"
)

// Endpoint containing the name and handler to call with the request
type Endpoint struct {
	// Name is the endpoint name, which should just be a single word, eg: "register"
	Name string
	// Mean is the mean average response time (time to generate response) promised for this endpoint
	Mean int32
	// Upper95 is 95th percentile response promised for this endpoint
	Upper95 int32
	// Handler is the function that will be fed requests to respond to
	Handler Handler
	// RequestProtocol is a struct type into which an inbound request for this endpoint can be unmarshaled
	RequestProtocol proto.Message
	// ResponseProtocol is the struct type defining the response format for this endpoint
	ResponseProtocol proto.Message
	// Subscribe indicates this endpoint should subscribe to a PUB stream - and gives us the stream address to SUB from
	//(topic name)
	Subscribe string
	// Authoriser is something that can check authorisation for this endpoint -- defaulting to ADMIN only (if nothing
	//specified by service)
	Authoriser Authoriser

	protoTMtx sync.RWMutex
	reqProtoT reflect.Type // cached type
	rspProtoT reflect.Type // cached type
}

func (ep *Endpoint) GetName() string {
	return ep.Name
}

func (ep *Endpoint) GetMean() int32 {
	return ep.Mean
}

func (ep *Endpoint) GetUpper95() int32 {
	return ep.Upper95
}

// instrumentedHandler wraps the handler to provide instrumentation
func (ep *Endpoint) instrumentedHandler(req *Request) (proto.Message, perrors.Error) {
	start := time.Now()

	var err perrors.Error
	var m proto.Message

	// In a defer in case the handler panics
	defer func() {
		stats.Record(ep, err, time.Since(start))
		if err == nil {
			inst.Timing(1.0, "success."+ep.Name, time.Since(start))
			return
		}
		switch err.Type() {
		case perrors.ErrorBadRequest, perrors.ErrorNotFound:
			// Ignore errors that are caused by clients
			// TODO: consider a new stat for clienterror?
			inst.Timing(1.0, "success."+ep.Name, time.Since(start))
			return
		default:
			inst.Timing(1.0, "error."+ep.Name, time.Since(start))
		}
	}()

	traceIn(req)
	// check auth, only call handler if auth passes
	err = ep.Authoriser.Authorise(req)
	if err == nil {
		req.Auth().SetAuthorised(true)
		m, err = ep.Handler(req)
	}
	traceOut(req, m, err, time.Since(start))
	return m, err
}

// ProtoTypes returns the Types of the registered request and response protocols
func (ep *Endpoint) ProtoTypes() (reflect.Type, reflect.Type) {
	ep.protoTMtx.RLock()
	reqT, rspT := ep.reqProtoT, ep.rspProtoT
	ep.protoTMtx.RUnlock()

	if (reqT == nil && ep.RequestProtocol != nil) || (rspT == nil && ep.ResponseProtocol != nil) {
		ep.protoTMtx.Lock()
		reqT, rspT = ep.reqProtoT, ep.rspProtoT // Prevent thundering herd
		if reqT == nil && ep.RequestProtocol != nil {
			reqT = reflect.TypeOf(ep.RequestProtocol)
			ep.reqProtoT = reqT
		}
		if rspT == nil && ep.ResponseProtocol != nil {
			rspT = reflect.TypeOf(ep.ResponseProtocol)
			ep.rspProtoT = rspT
		}
		ep.protoTMtx.Unlock()
	}

	return reqT, rspT
}

// unmarshalRequest reads a request's payload into a RequestProtocol object
func (ep *Endpoint) unmarshalRequest(req *Request) (proto.Message, perrors.Error) {
	reqProtoT, _ := ep.ProtoTypes()

	if reqProtoT == nil { // No registered protocol
		return nil, nil
	}

	result := reflect.New(reqProtoT.Elem()).Interface().(proto.Message)
	if err := req.Unmarshal(result); err != nil {
		return nil, perrors.InternalServerError(fmt.Sprintf("%s.%s.unmarshal", Name, ep.Name), err.Error())
	}

	return result, nil
}

// traceIn traces a request inbound to a service to handle
func traceIn(req *Request) {
	if req.shouldTrace() {
		trace.Send(&traceproto.Event{
			Timestamp:         proto.Int64(time.Now().UnixNano()),
			TraceId:           proto.String(req.TraceID()),
			Type:              traceproto.Event_IN.Enum(),
			MessageId:         proto.String(req.MessageID()),
			ParentMessageId:   proto.String(req.ParentMessageID()),
			From:              proto.String(req.From()),
			To:                proto.String(fmt.Sprintf("%v.%v", req.Service(), req.Endpoint())),
			Hostname:          proto.String(hostname),
			Az:                proto.String(az),
			Payload:           proto.String(""), // @todo
			HandlerInstanceId: proto.String(InstanceID),
			PersistentTrace:   proto.Bool(req.TraceShouldPersist()),
		})
	}
}

// traceOut traces a request outbound from a service handler
func traceOut(req *Request, msg proto.Message, err perrors.Error, d time.Duration) {
	if req.shouldTrace() {
		e := &traceproto.Event{
			Timestamp:         proto.Int64(time.Now().UnixNano()),
			TraceId:           proto.String(req.TraceID()),
			Type:              traceproto.Event_OUT.Enum(),
			MessageId:         proto.String(req.MessageID()),
			ParentMessageId:   proto.String(req.ParentMessageID()),
			From:              proto.String(req.From()),
			To:                proto.String(fmt.Sprintf("%v.%v", req.Service(), req.Endpoint())),
			Hostname:          proto.String(hostname),
			Az:                proto.String(az),
			Payload:           proto.String(""), // @todo
			HandlerInstanceId: proto.String(InstanceID),
			Duration:          proto.Int64(int64(d)),
			PersistentTrace:   proto.Bool(req.TraceShouldPersist()),
		}
		if err != nil {
			e.ErrorCode = proto.String(err.Code())
			e.ErrorDescription = proto.String(err.Description())
		}
		trace.Send(e)
	}
}
