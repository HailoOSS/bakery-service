package packer

import (
	"fmt"
	"io"

	// "github.com/hailocab/bakery-service/aws"

	// log "github.com/cihub/seelog"
	// "github.com/mitchellh/packer/packer"
	"github.com/mitchellh/packer/template"
)

const (
	PluginMaxPort = 10000
	PluginMinPort = 15000
)

type Packer struct {
	Template *template.Template
}

func New(t io.ReadCloser) (*Packer, error) {
	tpl, err := ReadTemplate(t)
	if err != nil {
		return nil, fmt.Errorf("Unable to create new packer: %v", err)
	}

	return &Packer{
		Template: tpl,
	}, nil
}

func (p *Packer) Build() error {
	config := NewConfig(PluginMinPort, PluginMaxPort)
	if err := config.Discover(); err != nil {
		return fmt.Errorf("Unable to discover packer config: %v", err)
	}

	return nil
}

func ReadTemplate(t io.ReadCloser) (*template.Template, error) {
	defer t.Close()

	tpl, err := template.Parse(t)
	if err != nil {
		return nil, fmt.Errorf("Unable to read template: %v", err)
	}

	if err := tpl.Validate(); err != nil {
		return nil, fmt.Errorf("The template is not valid: %v", err)
	}

	return tpl, nil
}
