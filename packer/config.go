package packer

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/mitchellh/osext"
	"github.com/mitchellh/packer/packer"
	"github.com/mitchellh/packer/packer/plugin"
)

type config struct {
	PluginMinPort uint
	PluginMaxPort uint

	Builders       map[string]string
	Provisioners   map[string]string
	PostProcessors map[string]string
}

func NewConfig(minPort uint, maxPort uint) *config {
	return &config{
		PluginMinPort: minPort,
		PluginMaxPort: maxPort,
	}
}

func (c *config) Discover() error {
	path, err := exec.LookPath("packer")
	if err != nil {
		return fmt.Errorf("Unable to find packer in the path: %v", err)
	}

	packerDir := filepath.Dir(path)
	if !filepath.IsAbs(packerDir) {
		path, err = filepath.Abs(packerDir)
		if err != nil {
			return fmt.Errorf("Packer path is invalid: %v", err)
		}
	}

	var globber Glob

	if err := c.discoverSingle(filepath.Join(path, "packer-builder-*"), &c.Builders, globber); err != nil {
		return fmt.Errorf("Couldn't discover builders: %v", err)
	}

	if err := c.discoverSingle(filepath.Join(path, "packer-post-processor-*"), &c.PostProcessors, globber); err != nil {
		return fmt.Errorf("Couldn't discover post processors: %v", err)
	}

	if err := c.discoverSingle(filepath.Join(path, "packer-provisioner-*"), &c.Provisioners, globber); err != nil {
		return fmt.Errorf("Couldn't discover provisioners: %v", err)
	}

	return nil
}

func (c *config) discoverSingle(pattern string, m *map[string]string, glob Globber) error {
	matches, err := glob.Glob(pattern)
	if err != nil {
		return fmt.Errorf("Unable to glob '%s': %v", pattern, err)
	}

	if *m == nil {
		*m = make(map[string]string)
	}

	prefix := filepath.Base(pattern)
	prefix = prefix[:strings.Index(prefix, "*")]
	for _, match := range matches {
		file := filepath.Base(match)
		if idx := strings.Index(file, "."); idx >= 0 {
			file = file[:idx]
		}

		plugin := file[len(prefix):]
		(*m)[plugin] = match
	}

	return nil
}

func (c *config) LoadBuilder(name string) (packer.Builder, error) {
	bin, ok := c.Builders[name]
	if !ok {
		return nil, fmt.Errorf("Unable to load builder: %s", name)
	}

	return c.pluginClient(bin).Builder()
}

func (c *config) LoadHook(name string) (packer.Hook, error) {
	return c.pluginClient(name).Hook()
}

func (c *config) LoadPostProcessor(name string) (packer.PostProcessor, error) {
	bin, ok := c.PostProcessors[name]
	if !ok {
		return nil, fmt.Errorf("Unable to load post processor: %s", name)
	}

	return c.pluginClient(bin).PostProcessor()
}

func (c *config) LoadProvisioner(name string) (packer.Provisioner, error) {
	bin, ok := c.Provisioners[name]
	if !ok {
		return nil, fmt.Errorf("Unable to load provisioner: %s", name)
	}

	return c.pluginClient(bin).Provisioner()
}

func (c *config) pluginClient(path string) *plugin.Client {
	originalPath := path

	// First attempt to find the executable by consulting the PATH.
	path, err := exec.LookPath(path)
	if err != nil {
		// If that doesn't work, look for it in the same directory
		// as the `packer` executable (us).
		exePath, err := osext.Executable()
		if err != nil {
		} else {
			path = filepath.Join(filepath.Dir(exePath), filepath.Base(originalPath))
		}
	}

	// If everything failed, just use the original path and let the error
	// bubble through.
	if path == "" {
		path = originalPath
	}

	var config plugin.ClientConfig

	config.Cmd = exec.Command(path)
	config.Managed = true
	config.MinPort = c.PluginMinPort
	config.MaxPort = c.PluginMaxPort

	return plugin.NewClient(&config)
}
