package packer

import (
	"testing"
)

func TestNewConfig(t *testing.T) {
	config := NewConfig(1, 2)

	if config.PluginMinPort != 1 || config.PluginMaxPort != 2 {
		t.Fatal("Unable to create a valid config")
	}
}

func GlobMock(pattern string) ([]string, error) {
	return []string{"/usr/local/bin/packer-builder-null"}, nil
}

func TestDiscoverySingle(t *testing.T) {
	config := NewConfig(1, 2)

	if err := config.discoverSingle("/usr/local/bin/packer-builder-*", &config.Builders, GlobMock); err != nil {
		t.Fatalf("Unable to discover builds: %v", err)
	}

	if len(config.Builders) == 0 {
		t.Fatal("No builders found")
	}
}
