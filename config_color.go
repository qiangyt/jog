package main

import (
	"github.com/gookit/color"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// ColorConfigT ...
type ColorConfigT struct {
	yaml.Unmarshaler
	yaml.Marshaler

	Label string
	Style color.Style
}

// ColorConfig ...
type ColorConfig = *ColorConfigT

// NewColorConfig ...
func NewColorConfig(label string) ColorConfig {
	var r ColorConfigT

	if err := yaml.UnmarshalStrict([]byte(label), &r); err != nil {
		panic(errors.Wrap(err, "failed to unmarshal color: "+label))
	}

	return &r
}

// UnmarshalYAML ...
func (me ColorConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	err := unmarshal(&me.Label)
	if err != nil {
		return err
	}

	me.Style, err = ColorsFromLabel(me.Label)
	return err
}

// MarshalYAML ...
func (me ColorConfig) MarshalYAML() (interface{}, error) {
	return me.Label, nil
}

// Render ...
func (me ColorConfig) Render(a ...interface{}) string {
	return me.Style.Render(a...)
}

// Sprint is alias of the 'Render'
func (me ColorConfig) Sprint(a ...interface{}) string {
	return me.Style.Sprint(a...)
}

// Sprintf format and render message.
func (me ColorConfig) Sprintf(format string, a ...interface{}) string {
	return me.Style.Sprintf(format, a...)
}

// Print render and Print text
func (me ColorConfig) Print(a ...interface{}) {
	me.Style.Print(a...)
}

// Printf render and print text
func (me ColorConfig) Printf(format string, a ...interface{}) {
	me.Style.Printf(format, a...)
}

// Println render and print text line
func (me ColorConfig) Println(a ...interface{}) {
	me.Style.Println(a...)
}

// OutputLevelsColorsConfigT ...
type OutputLevelsColorsConfigT struct {
	Debug ColorConfig
	Info  ColorConfig
	Error ColorConfig
	Warn  ColorConfig
	Trace ColorConfig
	Fine  ColorConfig
	Fatal ColorConfig
}

// OutputLevelsColorsConfig ...
type OutputLevelsColorsConfig = *OutputLevelsColorsConfigT

// OutputColorsConfigT ...
type OutputColorsConfigT struct {
	Index ColorConfig

	Timestamp   ColorConfig
	Version     ColorConfig
	Message     ColorConfig
	Logger      ColorConfig
	Thread      ColorConfig
	StackTrace  ColorConfig `yaml:"stack-trace"`
	StartedLine ColorConfig `yaml:"started-line"`

	PID    ColorConfig `yaml:"pid"`
	Host   ColorConfig
	File   ColorConfig
	Method ColorConfig
	Line   ColorConfig

	Levels OutputLevelsColorsConfig

	Raw             ColorConfig
	OthersName      ColorConfig `yaml:"others-name"`
	OthersSeparator ColorConfig `yaml:"others-separator"`
	OthersValue     ColorConfig `yaml:"others-value"`
}

// OutputColorsConfig ...
type OutputColorsConfig = *OutputColorsConfigT
