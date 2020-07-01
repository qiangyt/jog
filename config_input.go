package main

// FieldNamesConfigT ...
type FieldNamesConfigT struct {
	LineNo     string `yaml:"line-no"`
	Timestamp  string
	Version    string
	Message    string
	Logger     string
	Thread     string
	Level      string
	StackTrace string `yaml:"stack-trace"`
}

// FieldNamesConfig ...
type FieldNamesConfig = *FieldNamesConfigT

// InputConfigT ...
type InputConfigT struct {
	FieldNames FieldNamesConfig `yaml:"field-names"`
}

// InputConfig ...
type InputConfig = *InputConfigT
