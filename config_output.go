package main

// OutputConfigT ...
type OutputConfigT struct {
	Pattern            string
	CompressLoggerName bool `yaml:"compress-logger-name"`
	Colors             OutputColorsConfig
	StartedLine        string `yaml:"started-line"`
	StartedLineAppend  string `yaml:"started-line-append"`
}

// OutputConfig ...
type OutputConfig = *OutputConfigT
