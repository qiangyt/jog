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
	style, err := ColorsFromLabel(label)
	if err != nil {
		panic(errors.Wrap(err, ""))
	}

	return &ColorConfigT{
		Label: label,
		Style: style,
	}
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

// DefaultOutputLevelsColorsConfig ...
func DefaultOutputLevelsColorsConfig() OutputLevelsColorsConfig {
	return &OutputLevelsColorsConfigT{
		Debug: NewColorConfig("FgBlue,OpBold"),
		Info:  NewColorConfig("FgBlue,OpBold"),
		Error: NewColorConfig("FgRed,OpBold"),
		Warn:  NewColorConfig("FgYellow,OpBold"),
		Trace: NewColorConfig("FgBlue,OpBold"),
		Fine:  NewColorConfig("FgCyan,OpBold"),
		Fatal: NewColorConfig("FgRed,OpBold"),
	}
}

// OutputColorsConfigT ...
type OutputColorsConfigT struct {
	Index  ColorConfig
	Prefix ColorConfig

	App         ColorConfig
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

// DefaultOutputColorsConfig ...
func DefaultOutputColorsConfig() OutputColorsConfig {
	return &OutputColorsConfigT{
		Index:  NewColorConfig("FgDefault, OpBold"),
		Prefix: NewColorConfig("FgCyan"),

		App:         NewColorConfig("FgDefault"),
		Timestamp:   NewColorConfig("FgDefault"),
		Version:     NewColorConfig("FgDefault"),
		Message:     NewColorConfig("FgDefault"),
		Logger:      NewColorConfig("FgDefault"),
		Thread:      NewColorConfig("FgDefault"),
		StackTrace:  NewColorConfig("FgDefault"),
		StartedLine: NewColorConfig("FgGreen, OpBold"),

		PID:    NewColorConfig("FgDefault"),
		Host:   NewColorConfig("FgDefault"),
		File:   NewColorConfig("FgDefault"),
		Method: NewColorConfig("FgDefault"),
		Line:   NewColorConfig("FgDefault"),

		Levels: DefaultOutputLevelsColorsConfig(),

		Raw:             NewColorConfig("FgDefault"),
		OthersName:      NewColorConfig("FgDefault,OpBold"),
		OthersSeparator: NewColorConfig("FgDefault"),
		OthersValue:     NewColorConfig("FgDefault"),
	}
}
