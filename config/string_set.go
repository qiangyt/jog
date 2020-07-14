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
}

// StringSet ...
type StringSet = *StringSetT

func _extractKeys(m map[string]bool) []string {
	r := make([]string, len(m))
	n := 0
	for k := range m {
		r[n] = k
		n++
	}
	return r
}

// Parse ...
func (i StringSet) Parse(input interface{}) {
	i.Reset()

	switch input.(type) {
	case []string:
		{
			for _, v := range input.([]string) {
				v = strings.Trim(v, "\t\r\n ")
				i.ValueMap[v] = true
				if i.CaseSensitive == false {
					i.LowercasedValueMap[strings.ToLower(v)] = true
				}
			}
		}
	case string:
		{
			for _, v := range strutil.Split(input.(string), ",") {
				v = strings.Trim(v, "\t\r\n ")
				i.ValueMap[v] = true
				if i.CaseSensitive == false {
					i.LowercasedValueMap[strings.ToLower(v)] = true
				}
			}
		}
	default:
		panic(fmt.Errorf("not a string array: %v", input))
	}
}

// Reset ...
func (i StringSet) Reset() {
	i.ValueMap = make(map[string]bool)
	i.LowercasedValueMap = make(map[string]bool)
}

// Contains ...
func (i StringSet) Contains(v string) bool {
	r := i.ValueMap[v]
	if r {
		return true
	}
	if i.CaseSensitive == false {
		return i.LowercasedValueMap[v]
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
