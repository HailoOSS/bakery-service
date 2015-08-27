package packer

import (
	log "github.com/cihub/seelog"
)

type Ui struct{}

func (ui *Ui) Ask(prompt string) (string, error) {
	return "", nil
}

func (ui *Ui) Say(s string) {
	log.Info(s)
}

func (ui *Ui) Message(s string) {
	log.Info(s)
}

func (ui *Ui) Error(s string) {
	log.Info(s)
}

func (ui *Ui) Machine(s string, p ...string) {
	log.Infof(s, p)
}
