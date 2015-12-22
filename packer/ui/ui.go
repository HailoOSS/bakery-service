package ui

import (
	"fmt"

	log "github.com/cihub/seelog"
	"github.com/mitchellh/packer/packer"
)

// Caller foo
type Caller interface {
	Call(msg *Message)
}

// Callers foo
type Callers map[string]Caller

// CallerFunc foo
type CallerFunc func(c Callers)

// AddCaller foo
func AddCaller(name string, caller Caller) CallerFunc {
	return func(c Callers) {
		c[name] = caller
	}
}

type callerType int

func (ct callerType) String() string {
	return callerTypeDescs[ct]
}

const (
	callerTypeAsk callerType = 1
	callerTypeSay
	callerTypeMessage
	callerTypeError
	callerTypeMachine
)

var (
	callerTypeDescs = []string{
		"",
		"Ask",
		"Say",
		"Message",
		"Error",
		"Machine",
	}
)

// Message information
type Message struct {
	Type    callerType
	Message string
}

// UI struct
type UI struct {
	Callers Callers
}

// New creates a UI and passes the callers
func New(callers ...CallerFunc) packer.Ui {
	_callers := make(Callers)

	for _, c := range callers {
		c(_callers)
	}

	return &UI{
		Callers: _callers,
	}
}

// Ask a for information
func (ui *UI) Ask(prompt string) (string, error) {
	ui.call(callerTypeAsk, prompt)

	return "", fmt.Errorf("This isn't implemented")
}

// Say func
func (ui *UI) Say(message string) {
	ui.call(callerTypeSay, message)
}

// Message func
func (ui *UI) Message(message string) {
	ui.call(callerTypeMessage, message)
}

// Error func
func (ui *UI) Error(message string) {
	ui.call(callerTypeError, message)
}

// Machine func
func (ui *UI) Machine(t string, args ...string) {
	ui.call(callerTypeMachine, fmt.Sprintf("%s: %#v", t, args))
}

func (ui *UI) call(ct callerType, message string) {
	for n, c := range ui.Callers {
		log.Debugf("Calling %q: %s - %s", n, ct.String(), message)
		c.Call(&Message{
			Type:    ct,
			Message: message,
		})
	}
}
