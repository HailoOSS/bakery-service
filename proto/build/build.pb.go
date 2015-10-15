// Code generated by protoc-gen-go.
// source: github.com/hailocab/bakery-service/proto/build/build.proto
// DO NOT EDIT!

/*
Package com_hailocab_service_bakery_build is a generated protocol buffer package.

It is generated from these files:
	github.com/hailocab/bakery-service/proto/build/build.proto

It has these top-level messages:
	Request
	Response
	Variable
*/
package com_hailocab_service_bakery_build

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type Request struct {
	Template         *string     `protobuf:"bytes,1,req,name=template" json:"template,omitempty"`
	Variables        []*Variable `protobuf:"bytes,2,rep,name=variables" json:"variables,omitempty"`
	XXX_unrecognized []byte      `json:"-"`
}

func (m *Request) Reset()         { *m = Request{} }
func (m *Request) String() string { return proto.CompactTextString(m) }
func (*Request) ProtoMessage()    {}

func (m *Request) GetTemplate() string {
	if m != nil && m.Template != nil {
		return *m.Template
	}
	return ""
}

func (m *Request) GetVariables() []*Variable {
	if m != nil {
		return m.Variables
	}
	return nil
}

type Response struct {
	Id               *string `protobuf:"bytes,1,req,name=id" json:"id,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *Response) Reset()         { *m = Response{} }
func (m *Response) String() string { return proto.CompactTextString(m) }
func (*Response) ProtoMessage()    {}

func (m *Response) GetId() string {
	if m != nil && m.Id != nil {
		return *m.Id
	}
	return ""
}

type Variable struct {
	Key              *string `protobuf:"bytes,1,req,name=key" json:"key,omitempty"`
	Value            *string `protobuf:"bytes,2,req,name=value" json:"value,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *Variable) Reset()         { *m = Variable{} }
func (m *Variable) String() string { return proto.CompactTextString(m) }
func (*Variable) ProtoMessage()    {}

func (m *Variable) GetKey() string {
	if m != nil && m.Key != nil {
		return *m.Key
	}
	return ""
}

func (m *Variable) GetValue() string {
	if m != nil && m.Value != nil {
		return *m.Value
	}
	return ""
}