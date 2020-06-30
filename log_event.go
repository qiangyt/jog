package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
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
		log.Printf("line %d is blank\n", lineNo)
		return NewRawLogEvent(lineNo, raw)
	}

	all := make(map[string]interface{})
	if err := json.Unmarshal([]byte(line), &all); err != nil {
		log.Printf("failed to parse line %d: <%s>\n\treason %v\n", lineNo, raw, errors.Wrap(err, ""))
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

var _logstashMediator LogstashMediator

// ParseRawLine ...
func ParseRawLine(lineNo int, raw string) LogEvent {
	r := NewLogEvent(lineNo, raw)
	if !r.IsParsed() {
		return r
	}

	var mediator LogMediator

	mediator = _logstashMediator
	amountOfFieldsPopulated := mediator.Populate(r)
	if amountOfFieldsPopulated <= 0 {
		log.Printf("no fields populated. line %d: <%s>\n", lineNo, raw)
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
	log.Printf("file is opened: %s\n", localFilePath)
	defer f.Close()

	ProcessLinesWithReader(f)
}

// ProcessLinesWithReader ...
func ProcessLinesWithReader(reader io.Reader) {

	buf := bufio.NewReader(reader)

	for lineNo := 1; true; lineNo++ {
		raw, err := buf.ReadString('\n')

		if len(raw) != 0 {
			// trim the tail \n
			if raw[len(raw)-1] == '\n' {
				raw = raw[:len(raw)-1]
			}
		}
		if err != nil {
			if err == io.EOF {
				log.Printf("got EOF, line %d\n", lineNo)
				ProcessRawLine(lineNo, raw)
				return
			}
			panic(errors.Wrapf(err, "failed to read line %d", lineNo))
		}

		ProcessRawLine(lineNo, raw)
	}
}
