package conf

import _util "github.com/qiangyt/jog/pkg/util"

// EnumT ...
type EnumT struct {
	Name  string
	Alias _util.MultiString
	Color _util.Color
}

// Enum ...
type Enum = *EnumT

// TODO: remove this? UnmarshalYAML ...
func (i Enum) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return _util.DynObject4YAML(i, unmarshal)
}

// TODO: remove this? MarshalYAML ...
func (i Enum) MarshalYAML() (interface{}, error) {
	return _util.DynObject2YAML(i)
}

// TODO: remove this? Reset ...
func (i Enum) Reset() {
	i.Name = ""

	i.Alias = &_util.MultiStringT{}
	i.Alias.Reset()

	i.Color = &_util.ColorT{}
	i.Color.Reset()
}

// FromMap ...
func (i Enum) FromMap(m map[string]interface{}) error {
	aliasV := _util.ExtractFromMap(m, "alias")
	if aliasV != nil {
		if err := _util.UnmashalYAMLAgain(aliasV, &i.Alias); err != nil {
			return err
		}
	}

	colorV := _util.ExtractFromMap(m, "color")
	if colorV != nil {
		if err := _util.UnmashalYAMLAgain(colorV, &i.Color); err != nil {
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
