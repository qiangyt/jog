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

// DefaultOutputConfig ...
func DefaultOutputConfig() OutputConfig {
	return &OutputConfigT{
		Pattern:            "${prefix} | ${timestamp} ${level} <${thread}> ${logger}: ${message} ${others} ${stacktrace}",
		CompressLoggerName: true,
		Colors:             DefaultOutputColorsConfig(),
		StartedLine:        "Started Application in",
		StartedLineAppend:  "\n\n",
	}
}
