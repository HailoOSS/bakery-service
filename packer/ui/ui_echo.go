package ui

import (
	"fmt"
)

// EchoCaller a dummy caller
type EchoCaller struct{}

// Call does something with the message
func (ec *EchoCaller) Call(msg *Message) {
	fmt.Printf("%s: %s", msg.Type.String(), msg.Message)
}
