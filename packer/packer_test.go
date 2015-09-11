package packer

import (
	"bytes"
	"testing"
)

var (
	TemplateWithoutDefinedVars = `{
	"builders":[{
		"type": "amazon-ebs",
		"access_key": "{{user` + "`" + `aws_access_key_id` + "`" + `}}",
		"secret_key": "{{user ` + "`" + `aws_secret_key` + "`" + `}}",
    "token": "{{user ` + "`" + `aws_session_token` + "`" + `}}",
		"region": "us-east-1",
		"source_ami": "ami-d3e145b8",
		"instance_type": "t2.small",
		"ssh_username": "ubuntu",
		"ami_name": "packer-quick-start {{timestamp}}"
	}]
}`

	TemplateWithDefinedVars = `{
	"variables": {
		"aws_access_key_id": ""
	},
	"builders":[{
		"type": "amazon-ebs",
		"access_key": "{{user` + "`" + `aws_access_key_id` + "`" + `}}",
		"secret_key": "{{user ` + "`" + `aws_secret_key` + "`" + `}}",
    "token": "{{user ` + "`" + `aws_session_token` + "`" + `}}",
		"region": "us-east-1",
		"source_ami": "ami-d3e145b8",
		"instance_type": "t2.small",
		"ssh_username": "ubuntu",
		"ami_name": "packer-quick-start {{timestamp}}"
	}]
}`
)

func TestReadTemplateNoVars(t *testing.T) {
	tpl, err := ReadTemplate(NewMockReadCloser(TemplateWithoutDefinedVars))

	if err != nil {
		t.Fatalf("Problem reading template: %v", err)
	}

	if string(tpl.RawContents) != TemplateWithoutDefinedVars {
		t.Fatalf("Templates don't match")
	}
}

func TestReadTemplateVars(t *testing.T) {
	tpl, err := ReadTemplate(NewMockReadCloser(TemplateWithDefinedVars))

	if err != nil {
		t.Fatalf("Problem reading template: %v", err)
	}

	if string(tpl.RawContents) != TemplateWithDefinedVars {
		t.Fatalf("Templates don't match")
	}

	if len(tpl.Variables) == 0 {
		t.Fatalf("Template doesn't contain any variables")
	}
}

type MockReadCloser struct {
	Data *bytes.Buffer
}

func NewMockReadCloser(s string) *MockReadCloser {
	return &MockReadCloser{
		Data: bytes.NewBufferString(s),
	}
}

func (m *MockReadCloser) Read(p []byte) (n int, err error) {
	return m.Data.Read(p)
}

func (m *MockReadCloser) Close() error {
	return nil
}
