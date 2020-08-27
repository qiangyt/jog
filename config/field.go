package config

import (
	"fmt"
	"strings"

	"github.com/gookit/goutil/strutil"
	"github.com/qiangyt/jog/util"
)

// FieldType ...
type FieldType int

const (
	// FieldTypeAuto ...
	FieldTypeAuto FieldType = iota

	// FieldTypeTime ...
	FieldTypeTime
)

// FieldT ...
type FieldT struct {
	ElementT

	Name           string
	Alias          MultiString
	CaseSensitive  bool `yaml:"case-sensitive"`
	CompressPrefix `yaml:"compress-prefix"`
	Enums          EnumMap
	Type           FieldType
	TimeFormat     string `yaml:"time-format"`
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

	i.Type = FieldTypeAuto
	i.TimeFormat = ""
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

	if i.Type != FieldTypeAuto {
		if i.Type == FieldTypeTime {
			r["type"] = "time"
		}
	}

	if len(i.TimeFormat) > 0 {
		r["time-format"] = i.TimeFormat
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

	typeV := util.ExtractFromMap(m, "type")
	if typeV != nil {
		typeT := strutil.MustString(typeV)
		if typeT == "time" {
			i.Type = FieldTypeTime
		} else if typeT == "auto" {
			i.Type = FieldTypeAuto
		} else {
			return fmt.Errorf("unknown field type: %s", typeT)
		}
	}

	timeFormatV := util.ExtractFromMap(m, "time-format")
	if timeFormatV != nil {
		i.TimeFormat = strutil.MustString(timeFormatV)
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
func (i Field) PrintBody(color util.Color, builder *strings.Builder, body string) {
	if i.NotEnum() {
		if i.CompressPrefix.Enabled {
			body = i.CompressPrefix.Compress(body)
		}
	} else {
		body = i.Enums.GetEnum(body).Name
	}

	if color == nil {
		builder.WriteString(body)
	} else {
		builder.WriteString(color.Sprint(body))
	}
}
