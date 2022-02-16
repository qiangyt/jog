package conf

import (
	"fmt"
	"time"

	"github.com/araddon/dateparse"
	"github.com/gookit/goutil/strutil"
	_util "github.com/qiangyt/jog/pkg/util"
)

// FieldType ...
type FieldType int

const (
	// FieldType_Auto ...
	FieldType_Auto FieldType = iota

	// FieldType_Time ...
	FieldType_Time
)

// FieldT ...
type FieldT struct {
	ElementT

	Name           string
	Alias          _util.MultiString
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

	i.Alias = &_util.MultiStringT{}
	i.Alias.Set("")

	i.CaseSensitive = false

	i.CompressPrefix = &CompressPrefixT{}
	i.CompressPrefix.Reset()

	i.Enums = &EnumMapT{}
	i.Enums.Reset()

	i.Type = FieldType_Auto
	i.TimeFormat = ""
	i.Timezone = ""
	i.TimeLocation = nil
}

// UnmarshalYAML ...
func (i Field) UnmarshalYAML(unmarshal func(interface{}) error) error {
	e := _util.DynObject4YAML(i, unmarshal)
	return e
}

// MarshalYAML ...
func (i Field) MarshalYAML() (interface{}, error) {
	return _util.DynObject2YAML(i)
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

	if i.Type != FieldType_Auto {
		if i.Type == FieldType_Time {
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

	caseSensitiveV := _util.ExtractFromMap(m, "case-sensitive")
	if caseSensitiveV != nil {
		i.CaseSensitive = _util.ToBool(caseSensitiveV)
	}

	aliasV := _util.ExtractFromMap(m, "alias")
	if aliasV != nil {
		i.Alias.Set(strutil.MustString(aliasV))
	}

	compressPrefixV := _util.ExtractFromMap(m, "compress-prefix")
	if compressPrefixV != nil {
		if err := _util.UnmashalYAMLAgain(compressPrefixV, &i.CompressPrefix); err != nil {
			return err
		}
	}

	enumsV := _util.ExtractFromMap(m, "enums")
	if enumsV != nil {
		if err := _util.UnmashalYAMLAgain(enumsV, &i.Enums); err != nil {
			return err
		}
	}

	typeV := _util.ExtractFromMap(m, "type")
	if typeV != nil {
		typeT := strutil.MustString(typeV)
		if typeT == "time" {
			i.Type = FieldType_Time
		} else if typeT == "auto" {
			i.Type = FieldType_Auto
		} else {
			return fmt.Errorf("unknown field type: %s", typeT)
		}
	}

	timeFormatV := _util.ExtractFromMap(m, "time-format")
	if timeFormatV != nil {
		i.TimeFormat = strutil.MustString(timeFormatV)
	}

	timezoneV := _util.ExtractFromMap(m, "timezone")
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
func (i Field) GetColor(value string) _util.Color {
	if i.IsEnum() {
		return i.Enums.GetEnum(value).Color
	}
	return i.Color
}

// ParseTimestamp ...
func (i Field) ParseTimestamp(text string) time.Time {
	var timeValue time.Time
	var err error

	loc := i.TimeLocation
	tmFormat := i.TimeFormat

	if loc != nil {
		if len(tmFormat) != 0 {
			timeValue, err = time.ParseInLocation(tmFormat, text, loc)
			if err != nil {
				panic(fmt.Errorf("failed to parse time value: %s, with format: %s, loc: %v", text, tmFormat, loc))
			}
		} else {
			timeValue, err = dateparse.ParseIn(text, loc)
			if err != nil {
				panic(fmt.Errorf("failed to parse time value: %s, loc: %v", text, loc))
			}
		}
	} else {
		if len(tmFormat) != 0 {
			timeValue, err = time.Parse(tmFormat, text)
			if err != nil {
				panic(fmt.Errorf("failed to parse time value: %s, with format: %s", text, tmFormat))
			}
		} else {
			timeValue, err = dateparse.ParseAny(text)
			if err != nil {
				panic(fmt.Errorf("failed to parse time value: %s", text))
			}
		}
	}

	return timeValue
}
