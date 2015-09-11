package packer

import (
	"path/filepath"
)

type Globber func(pattern string) ([]string, error)

func Glob(pattern string) ([]string, error) {
	return filepath.Glob(pattern)
}
