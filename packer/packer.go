package packer

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sync"

	log "github.com/cihub/seelog"

	"github.com/mitchellh/packer/packer"
	"github.com/mitchellh/packer/template"
)

const (
	// PluginMaxPort max port for communication between plugins
	PluginMaxPort = 10000

	// PluginMinPort min port for communication between plugins
	PluginMinPort = 15000
)

// Packer data store
type Packer struct {
	Template *template.Template

	coreConfig *packer.CoreConfig
	ui         packer.Ui
}

// New creates a new packer object
func New(t io.ReadCloser, ui packer.Ui) (*Packer, error) {
	tpl, err := ReadTemplate(t)
	if err != nil {
		return nil, fmt.Errorf("Unable to create new packer: %v", err)
	}

	return &Packer{
		Template: tpl,
		ui:       ui,
	}, nil
}

// Build performs the final build
func (p *Packer) Build(variables map[string]*Variable) (map[string][]packer.Artifact, error) {
	config := NewConfig(PluginMinPort, PluginMaxPort)
	if err := config.Discover(); err != nil {
		return nil, fmt.Errorf("Unable to discover packer config: %v", err)
	}

	p.coreConfig = p.BuildCoreConfig(config, variables)

	core, err := packer.NewCore(p.coreConfig)
	if err != nil {
		return nil, fmt.Errorf("Unable to create new core: %v", err)
	}

	builds, err := p.ListBuilds(core)
	if err != nil {
		return nil, fmt.Errorf("Unable to list builds: %v", err)
	}

	artifacts, errs := p.ProcessBuilds(builds)
	if len(errs) > 0 {
		return nil, fmt.Errorf("Unable to process builds")
	}

	return artifacts, nil
}

// BuildCoreConfig compiles config
func (p *Packer) BuildCoreConfig(config *Config, vars map[string]*Variable) *packer.CoreConfig {
	return &packer.CoreConfig{
		Components: packer.ComponentFinder{
			Builder:       config.LoadBuilder,
			Hook:          config.LoadHook,
			PostProcessor: config.LoadPostProcessor,
			Provisioner:   config.LoadProvisioner,
		},
		Template:  p.Template,
		Variables: p.extractVariables(vars),
	}
}

// ListBuilds lists available builds from the template
func (p *Packer) ListBuilds(core *packer.Core) ([]packer.Build, error) {
	var builds []packer.Build
	for _, n := range core.BuildNames() {
		log.Debugf("Creating build for %q", n)

		b, err := core.Build(n)
		if err != nil {
			return nil, fmt.Errorf("Unable to create build %q: %v", n, err)
		}

		builds = append(builds, b)
	}

	return builds, nil
}

// ProcessBuilds builds individual builds
func (p *Packer) ProcessBuilds(builds []packer.Build) (map[string][]packer.Artifact, map[string]error) {
	artifacts := map[string][]packer.Artifact{}
	errors := map[string]error{}

	var wg sync.WaitGroup
	for _, b := range builds {
		log.Infof("Processing build %q", b.Name())
		wg.Add(1)

		cacheDir, _ := ioutil.TempDir("/tmp", "bakery")

		defer os.RemoveAll(cacheDir)

		log.Infof("Setting cache directory: %s", cacheDir)
		cache := &packer.FileCache{CacheDir: cacheDir}

		go func(b packer.Build) {
			defer wg.Done()

			log.Infof("Preparing build for %q", b.Name())
			warnings, err := b.Prepare()
			if err != nil {
				log.Errorf("Problem preparing the build for %q: %v", b.Name(), err)
				errors[b.Name()] = err
				return
			}

			for _, w := range warnings {
				log.Debugf("Warning for %q: %v", b.Name(), w)
			}

			runArtifacts, err := b.Run(p.ui, cache)

			if err != nil {
				log.Errorf("Build '%s' errored: %s", b.Name(), err)
				errors[b.Name()] = err
			} else {
				log.Infof("Build '%s' finished.", b.Name())
				artifacts[b.Name()] = runArtifacts
			}
		}(b)
	}

	log.Infof("Waiting for builds to complete")
	wg.Wait()

	if len(errors) > 0 {
		log.Error("There were some problems building")
		for n, e := range errors {
			log.Errorf("%s: %v", n, e)
		}

		return nil, errors
	}

	return artifacts, nil
}

// ListTemplateVariables extracts variables from a template
func (p *Packer) ListTemplateVariables() map[string]*Variable {
	_vars := map[string]*Variable{}
	for n, v := range p.Template.Variables {
		_vars[n] = &Variable{
			v, "",
		}
	}

	return _vars
}

func (p *Packer) extractVariables(vars map[string]*Variable) map[string]string {
	_vars := map[string]string{}
	for n, v := range vars {
		_vars[n] = v.Value
	}

	log.Infof("Extracted vars: %#v", _vars)

	return _vars
}

// ReadTemplate reads template from io.ReadCloser and validates it
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

// ExtractVariables maps template variables with the template
func ExtractVariables(vars map[string]*template.Variable, values map[string]string) map[string]*Variable {
	_vars := map[string]*Variable{}

	for k, v := range vars {
		for name, value := range values {
			_vars[k] = &Variable{
				Variable: v,
			}

			if name == k {
				_vars[k].Value = value
				_vars[k].Default = value
			}
		}
	}

	return _vars
}

// CheckVariables ensures variables are set
func CheckVariables(vars map[string]*Variable) (bool, error) {
	for n, v := range vars {
		if len(v.Value) == 0 {
			return false, fmt.Errorf("Variable %q is not set, but requried", n)
		}
	}

	return true, nil
}
