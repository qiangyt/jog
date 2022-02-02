package convert

import (
	"github.com/qiangyt/jog/util"
)

// EnumT ...
type EnumT struct {
	Name  string
	Alias util.MultiString
	Color util.Color
}

// Enum ...
type Enum = *EnumT

// TODO: remove this? UnmarshalYAML ...
func (i Enum) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return util.DynObject4YAML(i, unmarshal)
}

// TODO: remove this? MarshalYAML ...
func (i Enum) MarshalYAML() (interface{}, error) {
	return util.DynObject2YAML(i)
}

// TODO: remove this? Reset ...
func (i Enum) Reset() {
	i.Name = ""

	i.Alias = &util.MultiStringT{}
	i.Alias.Reset()

	i.Color = &util.ColorT{}
	i.Color.Reset()
}

// FromMap ...
func (i Enum) FromMap(m map[string]interface{}) error {
	aliasV := util.ExtractFromMap(m, "alias")
	if aliasV != nil {
		if err := util.UnmashalYAMLAgain(aliasV, &i.Alias); err != nil {
			return err
		}
	}

	colorV := util.ExtractFromMap(m, "color")
	if colorV != nil {
		if err := util.UnmashalYAMLAgain(colorV, &i.Color); err != nil {
			return err
		}
	}

	return nil
}

// ToMap ...
func (i Enum) ToMap() map[string]interface{} {
	r := make(map[string]interface{})
	r["alias"] = i.Alias.String()
	r["color"] = i.Color.String()
	return r
}
