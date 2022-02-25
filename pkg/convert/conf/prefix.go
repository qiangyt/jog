package conf

import _util "github.com/qiangyt/jog/pkg/util"

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
func (i Prefix) ToMap() map[string]interface{} {
	return i.ElementT.ToMap()
}

// UnmarshalYAML ...
func (i Prefix) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return _util.DynObject4YAML(i, unmarshal)
}

// MarshalYAML ...
func (i Prefix) MarshalYAML() (interface{}, error) {
	return _util.DynObject2YAML(i)
}