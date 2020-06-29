package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/pkg/errors"
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

	isParsed bool
}

// LogEvent ...
type LogEvent = *LogEventT

// Println ...
func (me LogEvent) Println() {
	if me.isParsed {
		fmt.Printf("%d %s %s %s %s %s %s \n", me.LineNo, me.Timestamp, me.Version, me.Message, me.Logger, me.Thread, me.Level)
	} else {
		fmt.Print(me.Raw)
	}
}

// IsParsed ...
func (me LogEvent) IsParsed() bool {
	return me.isParsed
}

// NewRawLogEvent ...
func NewRawLogEvent(lineNo int, raw string) LogEvent {
	return &LogEventT{
		LineNo:   lineNo,
		Others:   make(map[string]interface{}),
		All:      make(map[string]interface{}),
		Raw:      raw,
		isParsed: false,
	}
}

// NewLogEvent ...
func NewLogEvent(lineNo int, raw string) LogEvent {

	line := strings.TrimSpace(raw)
	if len(line) == 0 {
		return NewRawLogEvent(lineNo, raw)
	}

	all := make(map[string]interface{})
	if err := json.Unmarshal([]byte(line), &all); err != nil {
		// TODO: panic(errors.Wrapf(err, "failed to parse line: \n%s"+line))
		return NewRawLogEvent(lineNo, raw)
	}

	return &LogEventT{
		LineNo:   lineNo,
		Others:   make(map[string]interface{}),
		All:      all,
		Raw:      raw,
		isParsed: true,
	}
}

var _logstashParser LogstashParser

// ParseRawLine ...
func ParseRawLine(lineNo int, raw string) LogEvent {
	r := NewLogEvent(lineNo, raw)
	if !r.IsParsed() {
		return r
	}

	var parser LogParser

	parser = _logstashParser
	amountOfFieldsPopulated := parser.Parse(r)
	if amountOfFieldsPopulated <= 0 {
		return NewRawLogEvent(lineNo, raw)
	}

	return r
}

// ProcessRawLine ...
func ProcessRawLine(lineNo int, raw string) {
	event := ParseRawLine(lineNo, raw)
	event.Println()
}

// ProcessLinesWithLocalFile ...
func ProcessLinesWithLocalFile(localFilePath string) {

	f, err := os.Open(localFilePath)
	if err != nil {
		panic(errors.Wrapf(err, "failed to read file: %s", localFilePath))
	}
	defer f.Close()

	ProcessLinesWithReader(f)
}

// ProcessLinesWithReader ...
func ProcessLinesWithReader(reader io.Reader) {

	buf := bufio.NewReader(reader)

	for lineNo := 1; true; lineNo++ {
		raw, err := buf.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				ProcessRawLine(lineNo, raw)
				return
			}
			panic(errors.Wrapf(err, "failed to read line %d", lineNo))
		}

		ProcessRawLine(lineNo, raw)
	}
}
