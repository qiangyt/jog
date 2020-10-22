package main

import (
	"github.com/qiangyt/jog/config"
	"github.com/qiangyt/jog/util"
)

// FieldValueT ...
type FieldValueT struct {
	value     util.AnyValue
	enumValue config.Enum
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
func NewFieldValue(fieldConfig config.Field, value util.AnyValue) FieldValue {
	var enumValue config.Enum
	var output string

	if fieldConfig.IsEnum() {
		enumValue = fieldConfig.Enums.GetEnum(value.Text)
		output = enumValue.Name
	} else {
		if fieldConfig.CompressPrefix.Enabled {
			output = fieldConfig.CompressPrefix.Compress(value.Text)
		} else {
			output = value.Text
		}
	}

	return &FieldValueT{
		value:     value,
		enumValue: enumValue,
		Output:    output,
		Config:    fieldConfig,
	}
}
