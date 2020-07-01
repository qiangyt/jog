package main

// FieldNamesConfigT ...
type FieldNamesConfigT struct {
	Timestamp  string
	Version    string
	Message    string
	Logger     string
	Thread     string
	Level      string
	StackTrace string `yaml:"stack-trace"`
	PID        string `yaml:"pid"`
	Host       string
	File       string
	Method     string
	Line       string
}

// FieldNamesConfig ...
type FieldNamesConfig = *FieldNamesConfigT

// InputConfigT ...
type InputConfigT struct {
	FieldNames            FieldNamesConfig `yaml:"field-names"`
	IgnoreConversionError bool             `yaml:"ignore-conversion-error"`
}

// InputConfig ...
type InputConfig = *InputConfigT
