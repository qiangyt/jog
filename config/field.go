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
func (i Enum) ToMap(m map[string]interface{}) error {
	m["alias"] = i.Alias
	m["color"] = i.Color

	return nil
}

// EnumMapT ...
type EnumMapT struct {
	CaseSensitive bool `yaml:"case-sensitive"`
	Default       string
	Values        map[string]Enum

	ValueMap map[string]Enum
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
func (i EnumMap) ToMap(m map[string]interface{}) error {
	m["case-sensitive"] = i.CaseSensitive
	m["default"] = i.Default

	for k, v := range i.Values {
		m[k] = v
	}

	return nil
}

// FieldT ...
type FieldT struct {
	ElementT

	Name               string
	Alias              MultiString
	CaseSensitive      bool               `yaml:"case-sensitive"`
	LoggerNameCompress LoggerNameCompress `yaml:"logger-name-compress"`
	Enums              EnumMap
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

	i.LoggerNameCompress = &LoggerNameCompressT{}
	i.LoggerNameCompress.Reset()

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
func (i Field) ToMap(m map[string]interface{}) error {
	if err := i.ElementT.ToMap(m); err != nil {
		return err
	}

	m["case-sensitive"] = i.CaseSensitive
	m["alias"] = i.Alias
	m["logger-name-compress"] = i.LoggerNameCompress
	m["enums"] = i.Enums

	return nil
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

	loggerNameCompressV := util.ExtractFromMap(m, "logger-name-compress")
	if loggerNameCompressV != nil {
		i.LoggerNameCompress.Set(strutil.MustString(loggerNameCompressV))
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
	if !i.NotEnum() {
		a = i.Enums.GetEnum(a).Name
	}
	builder.WriteString(color.Sprint(a))
}
