// Code generated by protoc-gen-go.
// source: github.com/hailocab/go-login-service/proto/changepassword/changepassword.proto
// DO NOT EDIT!

/*
Package com_hailocab_service_login_changepassword is a generated protocol buffer package.

It is generated from these files:
	github.com/hailocab/go-login-service/proto/changepassword/changepassword.proto

It has these top-level messages:
	Request
	Response
*/
package com_hailocab_service_login_changepassword

import proto "github.com/hailocab/protobuf/proto"
import json "encoding/json"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = &json.SyntaxError{}
var _ = math.Inf

type Request struct {
	// The mechanism to invoke to change the password (e.g. h1, h2, etc.)
	Mech *string `protobuf:"bytes,1,req,name=mech" json:"mech,omitempty"`
	// defines the application this data logically belongs to
	Application *string `protobuf:"bytes,2,req,name=application" json:"application,omitempty"`
	// Who we are changing the password for
	Username *string `protobuf:"bytes,3,req,name=username" json:"username,omitempty"`
	// The new password
	NewPassword *string `protobuf:"bytes,4,req,name=newPassword" json:"newPassword,omitempty"`
	// The old password to validate
	OldPassword      *string `protobuf:"bytes,5,opt,name=oldPassword" json:"oldPassword,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *Request) Reset()         { *m = Request{} }
func (m *Request) String() string { return proto.CompactTextString(m) }
func (*Request) ProtoMessage()    {}

func (m *Request) GetMech() string {
	if m != nil && m.Mech != nil {
		return *m.Mech
	}
	return ""
}

func (m *Request) GetApplication() string {
	if m != nil && m.Application != nil {
		return *m.Application
	}
	return ""
}

func (m *Request) GetUsername() string {
	if m != nil && m.Username != nil {
		return *m.Username
	}
	return ""
}

func (m *Request) GetNewPassword() string {
	if m != nil && m.NewPassword != nil {
		return *m.NewPassword
	}
	return ""
}

func (m *Request) GetOldPassword() string {
	if m != nil && m.OldPassword != nil {
		return *m.OldPassword
	}
	return ""
}

// Response is empty if the call was successful
type Response struct {
	XXX_unrecognized []byte `json:"-"`
}

func (m *Response) Reset()         { *m = Response{} }
func (m *Response) String() string { return proto.CompactTextString(m) }
func (*Response) ProtoMessage()    {}

func init() {
}