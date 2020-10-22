package util

import (
	"strings"

	"github.com/gookit/goutil/strutil"
	"gopkg.in/yaml.v2"
)

// MultiStringT ...
type MultiStringT struct {
	yaml.Unmarshaler
	yaml.Marshaler

	Text             string
	Values           map[string]bool
	LowercasedValues map[string]bool
}

// MultiString ...
type MultiString = *MultiStringT

// Set ...
func (i MultiString) Set(valuesText string) {
	i.Text = valuesText

	i.Values = make(map[string]bool)
	i.LowercasedValues = make(map[string]bool)

	for _, v := range strutil.Split(valuesText, ",") {
		v = strings.Trim(v, "\t\r\n ")
		i.Values[v] = true
		i.LowercasedValues[strings.ToLower(v)] = true
	}
}

// Reset ...
func (i MultiString) Reset() {
	i.Set("")
}

// Contains ...
func (i MultiString) Contains(v string, caseSensitive bool) bool {
	r := i.Values[v]
	if r {
		return true
	}
	if !caseSensitive {
		return i.LowercasedValues[v]
	}
	return false
}

func (i MultiString) String() string {
	return i.Text
}

// UnmarshalYAML ...
func (i MultiString) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var valuesText string
	err := unmarshal(&valuesText)
	if err != nil {
		return err
	}

	i.Set(valuesText)

	return nil
}

// MarshalYAML ...
func (i MultiString) MarshalYAML() (interface{}, error) {
	return i.String(), nil
}
