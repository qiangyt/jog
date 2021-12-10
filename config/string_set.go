package config

import (
	"fmt"
	"strings"

	"github.com/gookit/goutil/strutil"
	"gopkg.in/yaml.v2"
)

// StringSetT ...
type StringSetT struct {
	yaml.Unmarshaler
	yaml.Marshaler

	CaseSensitive      bool
	ValueMap           map[string]bool
	LowercasedValueMap map[string]bool
	UppercasedValueMap map[string]bool
}

// StringSet ...
type StringSet = *StringSetT

// Parse ...
func (i StringSet) Parse(input interface{}) {
	i.Reset()

	var valuesArray []string

	switch input.(type) {
	case []string:
		{
			valuesArray = input.([]string)
		}
	case string:
		{
			valuesArray = strutil.Split(input.(string), ",")
		}
	default:
		panic(fmt.Errorf("not a string array: %v", input))
	}

	for _, v := range valuesArray {
		v = strings.Trim(v, "\t\r\n ")
		i.ValueMap[v] = true
		if i.CaseSensitive == false {
			i.LowercasedValueMap[strings.ToLower(v)] = true
			i.UppercasedValueMap[strings.ToUpper(v)] = true
		}
	}
}

// Reset ...
func (i StringSet) Reset() {
	i.ValueMap = make(map[string]bool)
	i.LowercasedValueMap = make(map[string]bool)
	i.UppercasedValueMap = make(map[string]bool)
}

func (i StringSet) IsEmpty() bool {
	return i.ValueMap == nil || len(i.ValueMap) == 0
}

// Contains ...
func (i StringSet) Contains(v string) bool {
	r := i.ValueMap[v]
	if r {
		return true
	}
	if i.CaseSensitive == false {
		return i.LowercasedValueMap[strings.ToLower(v)]
	}
	return false
}

// ContainsPrefixOf ...
func (i StringSet) ContainsPrefixOf(v string) bool {
	if len(i.ValueMap) == 0 {
		return false
	}

	for prefix := range i.ValueMap {
		if strings.HasPrefix(v, prefix) {
			return true
		}
	}

	if i.CaseSensitive == false {
		var lowercasedV = strings.ToLower(v)
		for prefix := range i.LowercasedValueMap {
			if strings.HasPrefix(lowercasedV, prefix) {
				return true
			}
		}
	}

	return false
}

func (i StringSet) String() string {
	buf := &strings.Builder{}

	first := true
	for v := range i.ValueMap {
		if first {
			first = false
		} else {
			buf.WriteString(", ")
		}
		buf.WriteString(v)
	}
	return buf.String()
}

// UnmarshalYAML ...
func (i StringSet) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var input interface{}
	err := unmarshal(&input)
	if err != nil {
		return err
	}

	i.Parse(input)

	return nil
}

// MarshalYAML ...
func (i StringSet) MarshalYAML() (interface{}, error) {
	return i.String(), nil
}
