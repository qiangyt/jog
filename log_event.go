package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

// LogEventT ...
type LogEventT struct {
	LineNo     int
	Timestamp  string
	Version    string
	Message    string
	Logger     string
	Thread     string
	Level      string
	StackTrace string

	Others map[string]interface{}
	All    map[string]interface{}
	Raw    string

	isBlank bool
}

// LogEvent ...
type LogEvent = *LogEventT

// Println ...
func (me LogEvent) Println(format LogFormat) {
	fmt.Printf("%d %s %s %s %s %s %s \n", me.LineNo, me.Timestamp, me.Version, me.Message, me.Logger, me.Thread, me.Level)
}

// IsBlank ...
func (me LogEvent) IsBlank() bool {
	return me.isBlank
}

// NewBlankLogEvent ...
func NewBlankLogEvent(lineNo int, raw string) LogEvent {
	return &LogEventT{
		LineNo:  lineNo,
		Others:  make(map[string]interface{}),
		All:     make(map[string]interface{}),
		Raw:     raw,
		isBlank: true,
	}
}

// NewLogEvent ...
func NewLogEvent(lineNo int, raw string) LogEvent {

	line := strings.TrimSpace(raw)
	if len(line) == 0 {
		return NewBlankLogEvent(lineNo, raw)
	}

	all := make(map[string]interface{})
	if err := json.Unmarshal([]byte(line), &all); err != nil {
		// TODO: panic(errors.Wrapf(err, "failed to parse line: \n%s"+line))
		return NewBlankLogEvent(lineNo, raw)
	}

	return &LogEventT{
		LineNo:  lineNo,
		Others:  make(map[string]interface{}),
		All:     all,
		Raw:     raw,
		isBlank: false,
	}
}

var _logstash LogstashFormat

// ParseRawLine ...
func ParseRawLine(lineNo int, raw string) LogEvent {
	r := NewLogEvent(lineNo, raw)
	if r.IsBlank() {
		return r
	}

	err := _logstash.Parse(r)
	if err != nil {
		// TODO: panic(errors.Wrapf(err, "failed to parse log line %d: %s", lineNo, raw))
		return NewBlankLogEvent(lineNo, raw)
	}

	return r
}

// ProcessRawLine ...
func ProcessRawLine(lineNo int, raw string) {
	event := ParseRawLine(lineNo, raw)
	event.Println(_logstash)
}
