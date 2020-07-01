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
	Index int

	Timestamp  string
	Version    string
	Message    string
	Logger     string
	Thread     string
	Level      string
	StackTrace string
	PID        string
	Host       string
	File       string
	Method     string
	Line       string

	Others map[string]interface{}
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
		"index":      colors.Index.Render(me.Index),
		"timestamp":  colors.Timestamp.Render(me.Timestamp),
		"level":      determineLevelColor(me.Level, colors.Levels).Render(me.Level),
		"thread":     colors.Thread.Render(me.Thread),
		"logger":     colors.Logger.Render(resolveLoggerName(&cfg.Output, me.Logger)),
		"message":    colors.Message.Render(msg),
		"others":     formatOthers(cfg, me.Others),
		"stacktrace": colors.StackTrace.Render(me.StackTrace),
		"pid":        colors.PID.Render(me.PID),
		"host":       colors.Host.Render(me.Host),
		"file":       colors.File.Render(me.File),
		"method":     colors.StackTrace.Render(me.Method),
		"line":       colors.Line.Render(me.Line),
	}

	fmt.Println(os.Expand(cfg.Output.Pattern, func(fieldName string) string {
		return fields[fieldName]
	}))

	if me.IsStartedLine {
		if len(cfg.Output.StartedLineAppend) > 0 {
			fmt.Printf(cfg.Output.StartedLineAppend)
		}
	}
}

func determineLevelColor(level string, levelsConfig OutputLevelsColorsConfig) ColorConfig {
	switch level {
	case "debug":
		return levelsConfig.Debug
	case "info":
		return levelsConfig.Info
	case "error":
		return levelsConfig.Error
	case "warn":
		return levelsConfig.Warn
	case "warning":
		return levelsConfig.Warn
	case "trace":
		return levelsConfig.Trace
	case "fine":
		return levelsConfig.Fine
	case "fatal":
		return levelsConfig.Fatal
	default:
		return levelsConfig.Warn
	}
}

// NewRawLogEvent ...
func NewRawLogEvent(cfg Config, index int, raw string) LogEvent {
	return &LogEventT{
		Index:         index,
		Others:        make(map[string]interface{}),
		Raw:           raw,
		IsParsed:      false,
		IsStartedLine: false,
	}
}

// NewLogEvent ...
func NewLogEvent(cfg Config, index int, raw string) (LogEvent, map[string]interface{}) {

	line := strings.TrimSpace(raw)
	if len(line) == 0 {
		log.Printf("line %d is blank\n", index)
		return NewRawLogEvent(cfg, index, raw), nil
	}

	fields := make(map[string]interface{})
	if err := json.Unmarshal([]byte(line), &fields); err != nil {
		log.Printf("failed to parse line %d: <%s>\n\treason %v\n", index, raw, errors.Wrap(err, ""))
		return NewRawLogEvent(cfg, index, raw), nil
	}

	return &LogEventT{
		Index:         index,
		Others:        make(map[string]interface{}),
		Raw:           raw,
		IsParsed:      true,
		IsStartedLine: len(cfg.Output.StartedLine) > 0 && strings.Contains(raw, cfg.Output.StartedLine),
	}, fields
}

var _mediator GenerialMediator

// ParseRawLine ...
func ParseRawLine(cfg Config, index int, raw string) LogEvent {
	r, fields := NewLogEvent(cfg, index, raw)
	if !r.IsParsed {
		return r
	}

	mediator := _mediator

	amountOfFieldsPopulated := mediator.PopulateFields(cfg, r, fields)
	if amountOfFieldsPopulated <= 0 {
		log.Printf("no fields populated. line %d: <%s>\n", index, raw)
		return NewRawLogEvent(cfg, index, raw)
	}

	return r
}

// ProcessRawLine ...
func ProcessRawLine(cfg Config, index int, raw string) {
	event := ParseRawLine(cfg, index, raw)
	event.Println(cfg)
}

// ProcessLinesWithLocalFile ...
func ProcessLinesWithLocalFile(cfg Config, localFilePath string) {
	f, err := os.Open(localFilePath)
	if err != nil {
		panic(errors.Wrap(err, ""))
	}
	log.Printf("file is opened: %s\n", localFilePath)
	defer f.Close()

	ProcessLinesWithReader(cfg, f)
}

// ProcessLinesWithReader ...
func ProcessLinesWithReader(cfg Config, reader io.Reader) {

	buf := bufio.NewReader(reader)

	for i := 1; true; i++ {
		raw, err := buf.ReadString('\n')

		if len(raw) != 0 {
			// trim the tail \n
			if raw[len(raw)-1] == '\n' {
				raw = raw[:len(raw)-1]
			}
		}
		if err != nil {
			if err == io.EOF {
				log.Printf("got EOF, line %d\n", i)
				ProcessRawLine(cfg, i, raw)
				return
			}
			panic(errors.Wrapf(err, "failed to read line %d", i))
		}

		ProcessRawLine(cfg, i, raw)
	}
}
