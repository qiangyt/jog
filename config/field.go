package config

import (
	"fmt"
	"strings"

	"github.com/gookit/goutil/strutil"
	"github.com/pkg/errors"
	"github.com/qiangyt/jog/util"
)

// EnumT ...
type EnumT struct {
	Name  string
	Alias MultiString
	Color util.Color
}

// Enum ...
type Enum = *EnumT

// UnmarshalYAML ...
func (i Enum) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return util.UnmarshalYAML(i, unmarshal)
}

// MarshalYAML ...
func (i Enum) MarshalYAML() (interface{}, error) {
	return util.MarshalYAML(i)
}

// Reset ...
func (i Enum) Reset() {
	i.Name = ""

	i.Alias = &MultiStringT{}
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

// EnumMapT ...
type EnumMapT struct {
	CaseSensitive bool `yaml:"case-sensitive"`
	Default       string
	Values        map[string]Enum
	ValueMap      map[string]Enum
}

// EnumMap ...
type EnumMap = *EnumMapT

// Reset ...
func (i EnumMap) Reset() {
	i.CaseSensitive = false
	i.Default = ""
	i.Values = make(map[string]Enum)
	i.ValueMap = make(map[string]Enum)
}

// UnmarshalYAML ...
func (i EnumMap) UnmarshalYAML(unmarshal func(interface{}) error) error {
	e := util.UnmarshalYAML(i, unmarshal)
	return e
}

// MarshalYAML ...
func (i EnumMap) MarshalYAML() (interface{}, error) {
	return util.MarshalYAML(i)
}

// IsEmpty ...
func (i EnumMap) IsEmpty() bool {
	return len(i.Values) == 0
}

// GetEnum ...
func (i EnumMap) GetEnum(value string) Enum {
	if !i.CaseSensitive {
		value = strings.ToLower(value)
	}

	r, has := i.ValueMap[value]
	if has {
		return r
	}
	return i.ValueMap[i.Default]
}

// FromMap ...
func (i EnumMap) FromMap(m map[string]interface{}) error {
	caseSensitiveV := util.ExtractFromMap(m, "case-sensitive")
	if caseSensitiveV != nil {
		i.CaseSensitive = util.ToBool(caseSensitiveV)
	}

	defaultV := util.ExtractFromMap(m, "default")
	if defaultV != nil {
		i.Default = strutil.MustString(defaultV)
	}

	for k, v := range m {
		var ev Enum
		if err := util.UnmashalYAMLAgain(v, &ev); err != nil {
			return err
		}

		ev.Name = k
		i.Values[k] = ev

		i.ValueMap[k] = ev
		if !i.CaseSensitive {
			i.ValueMap[strings.ToLower(k)] = ev
		}

		for aliasName := range ev.Alias.Values {
			i.ValueMap[aliasName] = ev
			if !i.CaseSensitive {
				i.ValueMap[strings.ToLower(aliasName)] = ev
			}
		}

		delete(m, k)

		if len(i.Default) == 0 {
			i.Default = k
		}
	}

	if len(i.Default) == 0 {
		return errors.New("default enum not specified")
	}
	if _, defaultMatches := i.ValueMap[i.Default]; !defaultMatches {
		return fmt.Errorf("invalid default enum name: %s", i.Default)
	}

	return nil
}

// ToMap ...
func (i EnumMap) ToMap() map[string]interface{} {
	r := make(map[string]interface{})
	r["case-sensitive"] = i.CaseSensitive
	r["default"] = i.Default

	for k, v := range i.Values {
		r[k] = v.ToMap()
	}

	return r
}

// FieldT ...
type FieldT struct {
	ElementT

	Name           string
	Alias          MultiString
	CaseSensitive  bool `yaml:"case-sensitive"`
	CompressPrefix `yaml:"compress-prefix"`
	Enums          EnumMap
}

// Field ...
type Field = *FieldT

// Reset ...
func (i Field) Reset() {
	i.ElementT.Reset()

	i.Name = ""

	i.Alias = &MultiStringT{}
	i.Alias.Set("")

	i.CaseSensitive = false

	i.CompressPrefix = &CompressPrefixT{}
	i.CompressPrefix.Reset()

	i.Enums = &EnumMapT{}
	i.Enums.Reset()
}

// UnmarshalYAML ...
func (i Field) UnmarshalYAML(unmarshal func(interface{}) error) error {
	e := util.UnmarshalYAML(i, unmarshal)
	return e
}

// MarshalYAML ...
func (i Field) MarshalYAML() (interface{}, error) {
	return util.MarshalYAML(i)
}

// NotEnum ...
func (i Field) NotEnum() bool {
	return i.Enums.IsEmpty()
}

// ToMap ...
func (i Field) ToMap() map[string]interface{} {
	r := i.ElementT.ToMap()

	r["case-sensitive"] = i.CaseSensitive
	r["alias"] = i.Alias.String()
	r["compress-prefix"] = i.CompressPrefix.ToMap()
	if !i.NotEnum() {
		r["enums"] = i.Enums.ToMap()
	}
	return r
}

// FromMap ...
func (i Field) FromMap(m map[string]interface{}) error {
	if err := i.ElementT.FromMap(m); err != nil {
		return err
	}

	caseSensitiveV := util.ExtractFromMap(m, "case-sensitive")
	if caseSensitiveV != nil {
		i.CaseSensitive = util.ToBool(caseSensitiveV)
	}

	aliasV := util.ExtractFromMap(m, "alias")
	if aliasV != nil {
		i.Alias.Set(strutil.MustString(aliasV))
	}

	compressPrefixV := util.ExtractFromMap(m, "compress-prefix")
	if compressPrefixV != nil {
		if err := util.UnmashalYAMLAgain(compressPrefixV, &i.CompressPrefix); err != nil {
			return err
		}
	}

	enumsV := util.ExtractFromMap(m, "enums")
	if enumsV != nil {
		if err := util.UnmashalYAMLAgain(enumsV, &i.Enums); err != nil {
			return err
		}
	}

	return nil
}

// GetColor ...
func (i Field) GetColor(value string) util.Color {
	if i.NotEnum() {
		return i.Color
	}
	return i.Enums.GetEnum(value).Color
}

// PrintBody ...
func (i Field) PrintBody(color util.Color, builder *strings.Builder, a string) {
	if i.NotEnum() {
		if i.CompressPrefix.Enabled {
			a = i.CompressPrefix.Compress(a)
		}
	} else {
		a = i.Enums.GetEnum(a).Name
	}

	if color == nil {
		builder.WriteString(a)
	} else {
		builder.WriteString(color.Sprint(a))
	}
}
