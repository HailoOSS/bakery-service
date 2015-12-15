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

// UI thing
type UI struct {
}

// Ask something
func (ui *UI) Ask(string) (string, error) {
	return "", nil
}

// Say something
func (ui *UI) Say(message string) {
	fmt.Println(message)
}

// Message something
func (ui *UI) Message(message string) {
	fmt.Println(message)
}

// Erorr something
func (ui *UI) Error(message string) {
	fmt.Println(message)
}

// Machine something
func (ui *UI) Machine(t string, args ...string) {
	fmt.Printf("%s %#v", t, args)
}
