package main

import (
	"fmt"
	"time"

	"github.com/araddon/dateparse"
	"github.com/qiangyt/jog/config"
	"github.com/qiangyt/jog/util"
)

// FieldValueT ...
type FieldValueT struct {
	value     util.AnyValue
	enumValue config.Enum
	timeValue time.Time
	Output    string
	Config    config.Field
}

// FieldValue ...
type FieldValue = *FieldValueT

// GetColor ...
func (i FieldValue) GetColor() util.Color {
	if i.enumValue != nil {
		return i.enumValue.Color
	}
	return i.Config.Color
}

// NewFieldValue ...
func NewFieldValue(cfg config.Configuration, options Options, fieldConfig config.Field, value util.AnyValue) FieldValue {
	var enumValue config.Enum
	var err error
	var output string

	text := value.Text

	if fieldConfig.IsEnum() {
		enumValue = fieldConfig.Enums.GetEnum(text)
		output = enumValue.Name
	} else {
		if fieldConfig.CompressPrefix.Enabled {
			output = fieldConfig.CompressPrefix.Compress(text)
		} else {
			output = text
		}
	}

	var timeValue time.Time
	if options.HasTimestampFilter() {
		if fieldConfig.Type == config.FieldTypeTime {
			loc := fieldConfig.TimeLocation
			tmFormat := fieldConfig.TimeFormat

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
		}
	}

	return &FieldValueT{
		value:     value,
		enumValue: enumValue,
		timeValue: timeValue,
		Output:    output,
		Config:    fieldConfig,
	}
}

// ParseTimestamp ...
func ParseTimestamp(fieldConfig config.Field, text string) time.Time {
	var timeValue time.Time
	var err error

	loc := fieldConfig.TimeLocation
	tmFormat := fieldConfig.TimeFormat

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
