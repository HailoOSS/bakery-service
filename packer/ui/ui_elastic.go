package ui

import (
	"bytes"
	"encoding/json"
	"time"

	"github.com/hailocab/bakery-service/elastic"
)

// ElasticCaller elastic search caller
type ElasticCaller struct {
	ID      string
	Elastic *elastic.Elastic
}

// NewElasticCaller creates a new elastic caller
func NewElasticCaller(id string, e *elastic.Elastic) *ElasticCaller {
	return &ElasticCaller{
		ID:      id,
		Elastic: e,
	}
}

// Call writes msg to elastic search
func (ec *ElasticCaller) Call(msg *Message) {
	msg.ID = ec.ID
	msg.Date = time.Now()

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(msg); err != nil {
		return
	}

	ec.Elastic.Write(buf.String())
}
