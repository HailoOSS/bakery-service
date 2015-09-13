package packer

import (
	"path/filepath"

	"github.com/mitchellh/packer/template"
)

type Globber func(pattern string) ([]string, error)

func Glob(pattern string) ([]string, error) {
	return filepath.Glob(pattern)
}

type Variable struct {
	*template.Variable

	Value string
}
