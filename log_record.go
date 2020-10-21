package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gookit/goutil/strutil"
	"github.com/pkg/errors"
	"github.com/qiangyt/jog/config"
	"github.com/qiangyt/jog/util"
)

// FieldValueT ...
type FieldValueT struct {
	Value  util.AnyValue
	Config config.Field
}

// FieldValue ...
type FieldValue = *FieldValueT

// LogRecordT ...
type LogRecordT struct {
	LineNo int

	Prefix         string
	StandardFields map[string]FieldValue
	UnknownFields  map[string]util.AnyValue
	Raw            string

	Unknown     bool
	StartupLine bool
}

// LogRecord .
type LogRecord = *LogRecordT

// PrintElement ...
func (i LogRecord) PrintElement(cfg config.Configuration, element util.Printable, builder *strings.Builder, a string) {
	if !element.IsEnabled() {
		return
	}

	var color util.Color
	if cfg.Colorization {
		if i.StartupLine {
			color = cfg.StartupLine.Color
		} else {
			color = element.GetColor(a)
		}
	} else {
		color = nil
	}

	element.PrintBefore(color, builder)
	element.PrintBody(color, builder, a)
	element.PrintAfter(color, builder)
}

// PopulateUnknownFields ...
func (i LogRecord) PopulateUnknownFields(cfg config.Configuration, unknownFields map[string]util.AnyValue, result map[string]string) {
	if len(unknownFields) == 0 {
		return
	}

	n := cfg.Fields.Unknown.Name
	s := cfg.Fields.Unknown.Separator
	v := cfg.Fields.Unknown.Value

	builder := &strings.Builder{}
	first := true
	for fName, fValue := range unknownFields {
		if !first {
			builder.WriteString(", ")
		}
		first = false

		i.PrintElement(cfg, n, builder, fName)
		i.PrintElement(cfg, s, builder, "=")
		i.PrintElement(cfg, v, builder, fValue.String())
	}

	result["unknown"] = builder.String()
}

// PopulateStandardFields ...
func (i LogRecord) PopulateStandardFields(cfg config.Configuration, standardFields map[string]FieldValue, result map[string]string) {
	if len(standardFields) == 0 {
		return
	}

	for _, f := range standardFields {
		builder := &strings.Builder{}
		i.PrintElement(cfg, f.Config, builder, f.Value.String())

		result[f.Config.Name] = builder.String()
	}
}

// AsFlatLine ...
func (i LogRecord) AsFlatLine(cfg config.Configuration) string {
	builder := &strings.Builder{}

	printStartLine := i.StartupLine && cfg.StartupLine.IsEnabled()

	var startupLineColor util.Color
	if cfg.Colorization {
		startupLineColor = cfg.StartupLine.GetColor("")
	} else {
		startupLineColor = nil
	}
	if printStartLine {
		cfg.StartupLine.PrintBefore(startupLineColor, builder)
	}

	i.PrintElement(cfg, cfg.LineNo, builder, fmt.Sprintf("%-6v ", i.LineNo))

	if i.Unknown {
		i.PrintElement(cfg, cfg.UnknownLine, builder, i.Raw)
	} else {
		if len(i.Prefix) > 0 {
			i.PrintElement(cfg, cfg.Prefix, builder, strutil.MustString(i.Prefix))
		}

		result := make(map[string]string)

		i.PopulateUnknownFields(cfg, i.UnknownFields, result)
		i.PopulateStandardFields(cfg, i.StandardFields, result)

		builder.WriteString(os.Expand(cfg.Pattern, func(fieldName string) string {
			return result[fieldName]
		}))
	}

	if printStartLine {
		cfg.StartupLine.PrintAfter(startupLineColor, builder)
	}

	return builder.String()
}

func isStartupLine(cfg config.Configuration, raw string) bool {
	contains := cfg.StartupLine.Contains
	return len(contains) > 0 && strings.Contains(raw, contains)
}

// ParseAsRecord ...
func ParseAsRecord(cfg config.Configuration, lineNo int, rawLine string) LogRecord {
	r := &LogRecordT{
		LineNo:         lineNo,
		UnknownFields:  make(map[string]util.AnyValue),
		StandardFields: make(map[string]FieldValue),
		Raw:            rawLine,
		Unknown:        true,
		StartupLine:    isStartupLine(cfg, rawLine),
	}

	line := strings.TrimSpace(rawLine)
	if len(line) == 0 {
		log.Printf("line %d is blank\n", lineNo)
		return r
	}

	posOfLeftBracket := strings.IndexByte(line, '{')
	if posOfLeftBracket < 0 {
		log.Printf("line %d is not JSON line: <%s>\n", lineNo, rawLine)
		return r
	}
	if posOfLeftBracket > 0 {
		r.Prefix = line[:posOfLeftBracket]
		line = line[posOfLeftBracket:]
	}

	allFields := make(map[string]interface{})
	if err := json.Unmarshal([]byte(line), &allFields); err != nil {
		log.Printf("parse round 1 failed: line %d: <%s>\n\treason %v\n", lineNo, line, errors.Wrap(err, ""))
		line = strings.ReplaceAll(line, "\\\"", "\"")
		if err := json.Unmarshal([]byte(line), &allFields); err != nil {
			log.Printf("parse round 2 failed: line %d: <%s>\n\treason %v\n", lineNo, line, errors.Wrap(err, ""))
			return r
		}
	}
	r.Unknown = false

	standardsFieldConfig := cfg.Fields.Standards
	for fName, fValue := range allFields {
		v := util.AnyValueFromRaw(lineNo, fValue, cfg.Replace)

		fConfig, contains := standardsFieldConfig[fName]
		if contains {
			f := &FieldValueT{Value: v, Config: fConfig}
			r.StandardFields[fName] = f
		} else {
			r.UnknownFields[fName] = v
		}
	}

	return r
}
