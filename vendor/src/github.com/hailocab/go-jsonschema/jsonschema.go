package jsonschema

import (
	"encoding/json"
	"errors"
)

const (
	TYPE_ARRAY   = `array`
	TYPE_BOOLEAN = `boolean`
	TYPE_INTEGER = `integer`
	TYPE_NUMBER  = `number`
	TYPE_NULL    = `null`
	TYPE_OBJECT  = `object`
	TYPE_STRING  = `string`
)

func New() *JsonSchema {
	return &JsonSchema{
		Properties:  make(map[string]*JsonSchema),
		Definitions: make(map[string]*JsonSchema),
	}
}

// JsonSchema represents a very generic json schema implementation to match the Hailo RPC Api.
// Currently it does not support all draft 4 features and more will be added as required
type JsonSchema struct {
	// Meta params
	Id          string `json:"id,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	// Used as "Name" tag
	Property string   `json:"-"`
	Type     string   `json:"type,omitempty"`
	Required []string `json:"required,omitempty"`
	// Reference url
	Ref string `json:"$ref,omitempty"`
	// Json reference
	Schema string `json:"$schema,omitempty"`
	// Hierarchy
	Enum                []string               `json:"enum,omitempty"`
	Parent              *JsonSchema            `json:"-"`
	Definitions         map[string]*JsonSchema `json:"definitions,omitempty"`
	DefinitionsChildren []*JsonSchema          `json:"-"`
	Items               *JsonSchema            `json:"items,omitempty"`
	ItemsChildren       []*JsonSchema          `json:"-"`
	PropertiesChildren  []*JsonSchema          `json:"-"`
	Properties          map[string]*JsonSchema `json:"properties,omitempty"`
}

func (s *JsonSchema) AddEnum(i interface{}) error {
	is, err := marshalToJsonString(i)
	if err != nil {
		return err
	}

	for _, v := range s.Enum {
		if v == *is {
			return errors.New("Enum value already exists")
		}
	}

	s.Enum = append(s.Enum, *is)

	return nil
}

func (s *JsonSchema) AddRequired(value string) error {
	for _, v := range s.Required {
		if v == value {
			return errors.New("Required value already exists")
		}
	}
	s.Required = append(s.Required, value)
	return nil
}

func (s *JsonSchema) AddProperty(pr *JsonSchema) error {
	if pr.Property == "" {
		return errors.New("The Property attribute of the child schema must be set")
	}
	s.Properties[pr.Property] = pr
	return nil
}

func (s *JsonSchema) AddPropertiesChild(child *JsonSchema) {
	s.PropertiesChildren = append(s.PropertiesChildren, child)
}

func (s *JsonSchema) AddDefinitionChild(child *JsonSchema) {
	s.DefinitionsChildren = append(s.DefinitionsChildren, child)
}

func (s *JsonSchema) AddItemsChild(child *JsonSchema) {
	s.ItemsChildren = append(s.ItemsChildren, child)
}

func marshalToJsonString(value interface{}) (*string, error) {
	mBytes, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}

	sBytes := string(mBytes)
	return &sBytes, nil
}

var JSON_TYPES []string

func init() {
	JSON_TYPES = []string{
		TYPE_ARRAY,
		TYPE_BOOLEAN,
		TYPE_INTEGER,
		TYPE_NUMBER,
		TYPE_NULL,
		TYPE_OBJECT,
		TYPE_STRING}
}
