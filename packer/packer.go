package packer

import (
	"fmt"
	"io"

	// "github.com/hailocab/bakery-service/aws"

	// log "github.com/cihub/seelog"
	"github.com/mitchellh/packer/packer"
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

	p.BuildCoreConfig(config, map[string]*Variable{})

	return nil
}

func (p *Packer) BuildCoreConfig(config *config, vars map[string]*Variable) *packer.CoreConfig {
	return &packer.CoreConfig{
		Components: packer.ComponentFinder{
			Builder:       config.LoadBuilder,
			Hook:          config.LoadHook,
			PostProcessor: config.LoadPostProcessor,
			Provisioner:   config.LoadProvisioner,
		},
		Template:  p.Template,
		Variables: map[string]string{
		// "AWS_ACCESS_KEY_ID":  credentials.AccessKeyID,
		// "AWS_SECRET_KEY":     credentials.SecretAccessKey,
		// "AWS_SECURITY_TOKEN": credentials.SessionToken,
		},
	}
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

func ExtractVariables(vars map[string]*template.Variable, values map[string]string) map[string]*Variable {
	_vars := map[string]*Variable{}

	for k, v := range vars {
		for name, value := range values {
			_vars[k] = &Variable{
				Variable: v,
			}

			if name == k {
				_vars[k].Value = value
			}
		}
	}

	return _vars
}

func CheckVariables(vars map[string]*Variable) (bool, error) {}
