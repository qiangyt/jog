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

	IsParsed      bool
	IsStartedLine bool
}

// LogEvent ...
type LogEvent = *LogEventT

var _loggerNameCache = make(map[string]string)

func resolveLoggerName(outputCfg OutputConfig, loggerName string) string {
	if !outputCfg.CompressLoggerName {
		return loggerName
	}

	if existingOne, ok := _loggerNameCache[loggerName]; ok {
		return existingOne
	}

	pkgList := strings.Split(loggerName, ".")
	indexLast := len(pkgList) - 1
	for index, pkg := range pkgList[:indexLast] {
		pkgList[index] = string([]byte(pkg)[0])
	}

	r := strings.Join(pkgList, ".")
	_loggerNameCache[loggerName] = outputCfg.Colors.Logger.Render(r)
	return r
}

var build strings.Builder

func formatOthers(cfg Config, others map[string]interface{}) string {
	if len(others) == 0 {
		return ""
	}

	//var builder strings.Builder
	colorName := cfg.Output.Colors.OthersName
	colorSeparator := cfg.Output.Colors.OthersSeparator
	colorValue := cfg.Output.Colors.OthersValue

	var othersList []string
	for otherFieldName, otherFieldValue := range others {
		pair := colorName.Render(otherFieldName) + colorSeparator.Render("=") + colorValue.Render(otherFieldValue)
		othersList = append(othersList, pair)
	}
	return "{" + strings.Join(othersList, ", ") + "}"
}

// Println ...
func (me LogEvent) Println(cfg Config) {
	colors := cfg.Output.Colors

	if !me.IsParsed {
		colors.Raw.Println(me.Raw)
		return
	}

	msg := me.Message
	if me.IsStartedLine {
		msg = colors.StartedLine.Render(msg)
	}

	stackTrace := me.StackTrace
	if len(stackTrace) > 0 {
		stackTrace += "Stack trace: \n" + stackTrace
	}

	fields := map[string]string{
		"timestamp":  colors.Timestamp.Render(me.Timestamp),
		"level":      determineLevelColor(me.Level, colors).Render(me.Level),
		"thread":     colors.Thread.Render(me.Thread),
		"logger":     colors.Logger.Render(resolveLoggerName(&cfg.Output, me.Logger)),
		"message":    colors.Message.Render(msg),
		"others":     formatOthers(cfg, me.Others),
		"stacktrace": colors.StackTrace.Render(me.StackTrace),
	}

	fmt.Println(os.Expand(cfg.Output.Pattern, func(fieldName string) string {
		return fields[fieldName]
	}))

	if me.IsStartedLine {
		fmt.Printf("\n\n\n\n")
	}
}

func determineLevelColor(level string, colorsConfig *OutputColorsConfigT) ColorConfig {
	switch level {
	case "debug":
		return colorsConfig.Debug
	case "info":
		return colorsConfig.Info
	case "error":
		return colorsConfig.Error
	case "warn":
		return colorsConfig.Warn
	case "warning":
		return colorsConfig.Warn
	case "trace":
		return colorsConfig.Trace
	case "fine":
		return colorsConfig.Fine
	case "fatal":
		return colorsConfig.Fatal
	default:
		return colorsConfig.Warn
	}
}

// NewRawLogEvent ...
func NewRawLogEvent(cfg Config, lineNo int, raw string) LogEvent {
	return &LogEventT{
		LineNo:        lineNo,
		Others:        make(map[string]interface{}),
		All:           make(map[string]interface{}),
		Raw:           raw,
		IsParsed:      false,
		IsStartedLine: false,
	}
}

// NewLogEvent ...
func NewLogEvent(cfg Config, lineNo int, raw string) LogEvent {

	line := strings.TrimSpace(raw)
	if len(line) == 0 {
		log.Printf("line %d is blank\n", lineNo)
		return NewRawLogEvent(cfg, lineNo, raw)
	}

	all := make(map[string]interface{})
	if err := json.Unmarshal([]byte(line), &all); err != nil {
		log.Printf("failed to parse line %d: <%s>\n\treason %v\n", lineNo, raw, errors.Wrap(err, ""))
		return NewRawLogEvent(cfg, lineNo, raw)
	}

	return &LogEventT{
		LineNo:        lineNo,
		Others:        make(map[string]interface{}),
		All:           all,
		Raw:           raw,
		IsParsed:      true,
		IsStartedLine: strings.Contains(raw, cfg.Output.StartedLine),
	}
}

var _logstashMediator LogstashMediator

// ParseRawLine ...
func ParseRawLine(cfg Config, lineNo int, raw string) LogEvent {
	r := NewLogEvent(cfg, lineNo, raw)
	if !r.IsParsed {
		return r
	}

	var mediator LogMediator

	mediator = _logstashMediator
	amountOfFieldsPopulated := mediator.PopulateFields(cfg, r)
	if amountOfFieldsPopulated <= 0 {
		log.Printf("no fields populated. line %d: <%s>\n", lineNo, raw)
		return NewRawLogEvent(cfg, lineNo, raw)
	}

	return r
}

// ProcessRawLine ...
func ProcessRawLine(cfg Config, lineNo int, raw string) {
	event := ParseRawLine(cfg, lineNo, raw)
	event.Println(cfg)
}

// ProcessLinesWithLocalFile ...
func ProcessLinesWithLocalFile(cfg Config, localFilePath string) {
	f, err := os.Open(localFilePath)
	if err != nil {
		panic(errors.Wrapf(err, "failed to read file: %s", localFilePath))
	}
	log.Printf("file is opened: %s\n", localFilePath)
	defer f.Close()

	ProcessLinesWithReader(cfg, f)
}

// ProcessLinesWithReader ...
func ProcessLinesWithReader(cfg Config, reader io.Reader) {

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
				ProcessRawLine(cfg, lineNo, raw)
				return
			}
			panic(errors.Wrapf(err, "failed to read line %d", lineNo))
		}

		ProcessRawLine(cfg, lineNo, raw)
	}
}
