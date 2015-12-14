package packer

import (
	"bytes"
	"testing"
)

var (
	TemplateWithoutDefinedVars = `{
	"builders":[{
		"type": "file",
		"Target": "/dev/null",
		"Content": "Hello"
	}]
}`

	TemplateWithDefinedVars = `{
	"variables": {
		"aws_access_key_id": ""
	},
	"builders":[{
		"type": "file",
		"Target": "/dev/null",
		"Content": "{{user ` + "`aws_access_key_id`" + `}}"
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

func TestBuild(t *testing.T) {
	p, err := New(NewMockReadCloser(TemplateWithDefinedVars))
	if err != nil {
		t.Fatalf("Unable to create new Packer: %v", err)
	}

	vars := ExtractVariables(p.Template.Variables, map[string]string{
		"aws_access_key_id": "AKI123456",
	})

	ok, err := CheckVariables(vars)
	if !ok || err != nil {
		t.Fatalf("Problem with variables: %v", err)
	}

	_, err = p.Build(vars)
	if err != nil {
		t.Fatalf("Unable to build: %v", err)
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
