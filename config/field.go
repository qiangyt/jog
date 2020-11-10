package config

import (
	"fmt"
	"strings"
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
	Width          int
	Fixed          bool
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

	i.Width = 0
	i.Fixed = false
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

	if i.Width != 0 {
		r["width"] = i.Width
	}
	r["fixed"] = i.Fixed

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

	widthV := util.ExtractFromMap(m, "width")
	if widthV != nil {
		width := strutil.MustInt(strutil.MustString(widthV))
		if width > 0 && width <= 4 {
			i.Width = 4
		} else if width < 0 && width >= -4 {
			i.Width = -4
		} else {
			i.Width = width
		}
	}

	fixedV := util.ExtractFromMap(m, "fixed")
	if fixedV != nil {
		i.Fixed = util.ToBool(fixedV)
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

// PrintBody ...
func (i Field) PrintBody(color util.Color, builder *strings.Builder, body string) {
	if i.Fixed {
		body = ShortenValue(body, abs(i.Width))
	}

	format := GetFormatString(i.Width)

	if color == nil {
		builder.WriteString(fmt.Sprintf(format, body))
	} else {
		builder.WriteString(color.Sprintf(format, body))
	}
}

// ShortenValue shortens the value to maxWidth -3 chars if necessary, shortend values will be postfixed by three dots
func ShortenValue(inValue string, maxWidth int) string {
	if maxWidth == 0 || len([]rune(inValue)) <= maxWidth {
		return inValue
	}
	return fmt.Sprint(inValue[:(maxWidth-3)], "...")
}

// abs function that works for int, Math.Abs only accepts float64
func abs(value int) int {
	if value < 0 {
		value = value * -1
	}
	return value
}

// GetFormatString ...
func GetFormatString(width int) string {
	format := "%s"
	if width != 0 {
		format = fmt.Sprint("%", width, "s")
	}
	return format
}
