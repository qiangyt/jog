package conf

import (
	"github.com/gookit/goutil/strutil"
	_util "github.com/qiangyt/jog/pkg/util"
)

// StartupLineT ...
type StartupLineT struct {
	ElementT

	Contains string
}

// StartupLine ...
type StartupLine = *StartupLineT

// UnmarshalYAML ...
func (i StartupLine) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return _util.DynObject4YAML(i, unmarshal)
}

// MarshalYAML ...
func (i StartupLine) MarshalYAML() (interface{}, error) {
	return _util.DynObject2YAML(i)
}

// Reset ...
func (i StartupLine) Reset() {
	i.ElementT.Reset()

	i.Contains = "Started Application in"
}

// FromMap ...
func (i StartupLine) FromMap(m map[string]interface{}) error {
	if err := i.ElementT.FromMap(m); err != nil {
		return err
	}

	containsV := _util.ExtractFromMap(m, "contains")
	if containsV != nil {
		i.Contains = strutil.MustString(containsV)
	}

	return nil
}

// ToMap ...
func (i StartupLine) ToMap() map[string]interface{} {
	r := i.ElementT.ToMap()
	r["contains"] = i.Contains
	return r
}
