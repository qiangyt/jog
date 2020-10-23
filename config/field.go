package config

import (
	"fmt"
	"time"

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
	Alias          util.MultiString
	CaseSensitive  bool `yaml:"case-sensitive"`
	CompressPrefix `yaml:"compress-prefix"`
	Enums          EnumMap
	Type           FieldType
	TimeFormat     string `yaml:"time-format"`
	Timezone       string `yaml:"timezone"`
	TimeLocation   *time.Location
}

// Field ...
type Field = *FieldT

// Reset ...
func (i Field) Reset() {
	i.ElementT.Reset()

	i.Name = ""

	i.Alias = &util.MultiStringT{}
	i.Alias.Set("")

	i.CaseSensitive = false

	i.CompressPrefix = &CompressPrefixT{}
	i.CompressPrefix.Reset()

	i.Enums = &EnumMapT{}
	i.Enums.Reset()

	i.Type = FieldTypeAuto
	i.TimeFormat = ""
	i.Timezone = ""
	i.TimeLocation = nil
}

// UnmarshalYAML ...
func (i Field) UnmarshalYAML(unmarshal func(interface{}) error) error {
	e := UnmarshalYAML(i, unmarshal)
	return e
}

// MarshalYAML ...
func (i Field) MarshalYAML() (interface{}, error) {
	return MarshalYAML(i)
}

// Init ...
func (i Field) Init(cfg Configuration) {
}

// IsEnum ...
func (i Field) IsEnum() bool {
	return !i.Enums.IsEmpty()
}

// ToMap ...
func (i Field) ToMap() map[string]interface{} {
	r := i.ElementT.ToMap()

	r["case-sensitive"] = i.CaseSensitive
	r["alias"] = i.Alias.String()
	r["compress-prefix"] = i.CompressPrefix.ToMap()
	if i.IsEnum() {
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

	if len(i.Timezone) > 0 {
		r["timezone"] = i.Timezone
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

	timezoneV := util.ExtractFromMap(m, "timezone")
	if timezoneV != nil {
		i.Timezone = strutil.MustString(timezoneV)

		loc, err := time.LoadLocation(i.Timezone)
		if err != nil {
			return fmt.Errorf("invalid timezone: %s", i.Timezone)
		}
		i.TimeLocation = loc
	}

	return nil
}

// GetColor ...
func (i Field) GetColor(value string) util.Color {
	if i.IsEnum() {
		return i.Enums.GetEnum(value).Color
	}
	return i.Color
}
