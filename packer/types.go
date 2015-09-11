package packer

import (
	"path/filepath"
)

type Globber interface {
	Glob(pattern string) ([]string, error)
}

type Glob struct{}

func (g Glob) Glob(pattern string) ([]string, error) {
	return filepath.Glob(pattern)
}
