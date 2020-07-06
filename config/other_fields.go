package config

import (
	"github.com/gookit/goutil/strutil"
	"github.com/qiangyt/jog/util"
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
	return util.UnmarshalYAML(i, unmarshal)
}

// MarshalYAML ...
func (i SeparatorField) MarshalYAML() (interface{}, error) {
	return util.MarshalYAML(i)
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

	labelV := util.ExtractFromMap(m, "label")
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

// OtherFieldsT ...
type OtherFieldsT struct {
	Name      Element
	Separator SeparatorField
	Value     Element
}

// OtherFields ...
type OtherFields = *OtherFieldsT

// Reset ...
func (i OtherFields) Reset() {
	i.Name = &ElementT{}
	i.Name.Reset()

	i.Separator = &SeparatorFieldT{}
	i.Separator.Reset()

	i.Value = &ElementT{}
	i.Value.Reset()
}

// ToMap ...
func (i OtherFields) ToMap() map[string]interface{} {
	r := make(map[string]interface{})
	r["name"] = i.Name.ToMap()
	r["separator"] = i.Separator.ToMap()
	r["value"] = i.Value.ToMap()
	return r
}
