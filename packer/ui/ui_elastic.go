package ui

import (
	"bytes"
	"encoding/json"

	"github.com/hailocab/bakery-service/elastic"
)

// ElasticCaller elastic search caller
type ElasticCaller struct {
	Elastic *elastic.Elastic
}

// NewElasticCaller creates a new elastic caller
func NewElasticCaller(e *elastic.Elastic) *ElasticCaller {
	return &ElasticCaller{
		Elastic: e,
	}
}

// Call writes msg to elastic search
func (ec *ElasticCaller) Call(msg *Message) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(msg); err != nil {
		return
	}

	ec.Elastic.Write(buf.String())
}
