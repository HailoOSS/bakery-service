package packer

import (
	"fmt"
	"path/filepath"

	"github.com/mitchellh/packer/template"
)

// Globber type
type Globber func(pattern string) ([]string, error)

// Glob files on a filesystem
func Glob(pattern string) ([]string, error) {
	return filepath.Glob(pattern)
}

// Variable data
type Variable struct {
	*template.Variable

	Value string
}
