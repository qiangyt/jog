package main

import (
	"strings"

	"github.com/gookit/goutil/strutil"
	"gopkg.in/yaml.v2"
)

// FieldNameT ...
type FieldNameT struct {
	yaml.Unmarshaler
	yaml.Marshaler

	Names []string
}

// FieldName ...
type FieldName = *FieldNameT

// UnmarshalYAML ...
func (me FieldName) UnmarshalYAML(unmarshal func(interface{}) error) error {
	namesText := ""
	err := unmarshal(&namesText)
	if err != nil {
		return err
	}

	me.Names = strutil.Split(namesText, ",")

	return nil
}

// MarshalYAML ...
func (me FieldName) MarshalYAML() (interface{}, error) {
	return strings.Join(me.Names, ", "), nil
}

// FieldNamesConfigT ...
type FieldNamesConfigT struct {
	Timestamp  FieldName
	Version    FieldName
	Message    FieldName
	Logger     FieldName
	Thread     FieldName
	Level      FieldName
	StackTrace FieldName `yaml:"stack-trace"`
	PID        FieldName `yaml:"pid"`
	Host       FieldName
	File       FieldName
	Method     FieldName
	Line       FieldName
}

// FieldNamesConfig ...
type FieldNamesConfig = *FieldNamesConfigT

// InputConfigT ...
type InputConfigT struct {
	FieldNames            FieldNamesConfig `yaml:"field-names"`
	IgnoreConversionError bool             `yaml:"ignore-conversion-error"`
}

// InputConfig ...
type InputConfig = *InputConfigT
