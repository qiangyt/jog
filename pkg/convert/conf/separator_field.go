package conf

import (
	"github.com/gookit/goutil/strutil"
	_util "github.com/qiangyt/jog/pkg/util"
)

// SeparatorFieldT ...
type SeparatorFieldT struct {
	ElementT

	Label string
}

// SeparatorField ...
type SeparatorField = *SeparatorFieldT

// UnmarshalYAML ...
func (i SeparatorField) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return _util.DynObject4YAML(i, unmarshal)
}

// MarshalYAML ...
func (i SeparatorField) MarshalYAML() (interface{}, error) {
	return _util.DynObject2YAML(i)
}

// Reset ...
func (i SeparatorField) Reset() {
	i.ElementT.Reset()

	i.Label = "="
}

// FromMap ...
func (i SeparatorField) FromMap(m map[string]interface{}) error {
	if err := i.ElementT.FromMap(m); err != nil {
		return err
	}

	labelV := _util.ExtractFromMap(m, "label")
	if labelV != nil {
		i.Label = strutil.MustString(labelV)
	}
	return nil
}

// ToMap ...
func (i SeparatorField) ToMap() map[string]interface{} {
	r := i.ElementT.ToMap()
	r["label"] = i.Label
	return r
}