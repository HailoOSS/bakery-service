// Code generated by protoc-gen-go.
// source: github.com/hailocab/go-login-service/proto/sessionread/sessionread.proto
// DO NOT EDIT!

/*
Package com_hailocab_service_login_sessionread is a generated protocol buffer package.

It is generated from these files:
	github.com/hailocab/go-login-service/proto/sessionread/sessionread.proto

It has these top-level messages:
	Request
	Response
*/
package com_hailocab_service_login_sessionread

import proto "github.com/hailocab/protobuf/proto"
import json "encoding/json"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = &json.SyntaxError{}
var _ = math.Inf

type Request struct {
	SessId           *string `protobuf:"bytes,1,req,name=sessId" json:"sessId,omitempty"`
	NoRenew          *bool   `protobuf:"varint,2,opt,name=noRenew" json:"noRenew,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *Request) Reset()         { *m = Request{} }
func (m *Request) String() string { return proto.CompactTextString(m) }
func (*Request) ProtoMessage()    {}

func (m *Request) GetSessId() string {
	if m != nil && m.SessId != nil {
		return *m.SessId
	}
	return ""
}

func (m *Request) GetNoRenew() bool {
	if m != nil && m.NoRenew != nil {
		return *m.NoRenew
	}
	return false
}

type Response struct {
	SessId           *string `protobuf:"bytes,1,req,name=sessId" json:"sessId,omitempty"`
	Token            *string `protobuf:"bytes,2,req,name=token" json:"token,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *Response) Reset()         { *m = Response{} }
func (m *Response) String() string { return proto.CompactTextString(m) }
func (*Response) ProtoMessage()    {}

func (m *Response) GetSessId() string {
	if m != nil && m.SessId != nil {
		return *m.SessId
	}
	return ""
}

func (m *Response) GetToken() string {
	if m != nil && m.Token != nil {
		return *m.Token
	}
	return ""
}

func init() {
}
