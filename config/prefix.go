package config

import (
	"github.com/qiangyt/jog/util"
)

// PrefixT ...
type PrefixT struct {
	ElementT
}

// Prefix ...
type Prefix = *PrefixT

// Reset ...
func (i Prefix) Reset() {
	i.ElementT.Reset()

	i.Color.Set("FgBlue")
}

// FromMap ...
func (i Prefix) FromMap(m map[string]interface{}) error {
	return i.ElementT.FromMap(m)
}

// ToMap ...
func (i Prefix) ToMap(m map[string]interface{}) error {
	return i.ElementT.ToMap(m)
}

// UnmarshalYAML ...
func (i Prefix) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return util.UnmarshalYAML(i, unmarshal)
}

// MarshalYAML ...
func (i Prefix) MarshalYAML() (interface{}, error) {
	return util.MarshalYAML(i)
}
