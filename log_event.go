package main

import (
	"encoding/json"
	"strings"

	"github.com/pkg/errors"
)

// LogEventT ...
type LogEventT struct {
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

func (me LogEvent) IsBlank() bool {
	return me.isBlank
}

// NewLogEvent ...
func NewLogEvent(raw string) LogEvent {

	line := strings.TrimSpace(raw)
	if len(line) == 0 {
		return &LogEventT{
			Others:  make(map[string]interface{}),
			All:     make(map[string]interface{}),
			Raw:     raw,
			isBlank: true,
		}
	}

	all := make(map[string]interface{})
	if err := json.Unmarshal([]byte(line), &all); err != nil {
		panic(errors.Wrapf(err, "failed to parse line: \n%s"+line))
	}

	return &LogEventT{
		Others:  make(map[string]interface{}),
		All:     all,
		Raw:     raw,
		isBlank: false,
	}
}
