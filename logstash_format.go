package main

// LogstashFormatT implements LogFormat interface
type LogstashFormatT struct {
}

// LogstashFormat is pointer of LogstashFormatT
type LogstashFormat = *LogstashFormatT

// Parse the log event
func (me LogstashFormat) Parse(raw string) LogEvent {
	return nil
}
