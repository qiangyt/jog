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

// StandardFieldT ...
type StandardFieldT struct {
	Value  util.AnyValue
	Config config.Field
}

// StandardField ...
type StandardField = *StandardFieldT

// LogRecordT ...
type LogRecordT struct {
	LineNo int

	Prefix         string
	StandardFields []StandardField
	OtherFields    map[string]util.AnyValue
	Raw            string

	Unknown     bool
	StartupLine bool
}

// LogRecord .
type LogRecord = *LogRecordT

// PrintElement ...
func (i LogRecord) PrintElement(config Config, element util.Printable, builder *strings.Builder, a string) {
	if !element.IsEnabled() {
		return
	}

	var color util.Color
	if config.Colorization {
		if i.StartupLine {
			color = config.StartupLine.Color
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

// PopulateOtherFields ...
func (i LogRecord) PopulateOtherFields(cfg Config, result map[string]string) {
	if len(i.OtherFields) == 0 {
		return
	}

	n := cfg.Fields.Others.Name
	s := cfg.Fields.Others.Separator
	v := cfg.Fields.Others.Value

	builder := &strings.Builder{}
	first := true
	for fName, fValue := range i.OtherFields {
		if !first {
			builder.WriteString(", ")
		}
		first = false

		i.PrintElement(cfg, n, builder, fName)
		i.PrintElement(cfg, s, builder, "=")
		i.PrintElement(cfg, v, builder, fValue.String())
	}

	result["others"] = builder.String()
}

// PopulateStandardFields ...
func (i LogRecord) PopulateStandardFields(cfg Config, result map[string]string) {
	if len(i.StandardFields) == 0 {
		return
	}

	for _, f := range i.StandardFields {
		builder := &strings.Builder{}
		i.PrintElement(cfg, f.Config, builder, f.Value.String())

		result[f.Config.Name] = builder.String()
	}
}

// AsFlatLine ...
func (i LogRecord) AsFlatLine(cfg Config) string {
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

		i.PopulateOtherFields(cfg, result)
		i.PopulateStandardFields(cfg, result)

		builder.WriteString(os.Expand(cfg.Pattern, func(fieldName string) string {
			return result[fieldName]
		}))
	}

	if printStartLine {
		cfg.StartupLine.PrintAfter(startupLineColor, builder)
	}

	return builder.String()
}

func isStartupLine(cfg Config, raw string) bool {
	contains := cfg.StartupLine.Contains
	return len(contains) > 0 && strings.Contains(raw, contains)
}

// ParseAsRecord ...
func ParseAsRecord(cfg Config, lineNo int, rawLine string) LogRecord {
	r := &LogRecordT{
		LineNo:         lineNo,
		OtherFields:    make(map[string]util.AnyValue),
		StandardFields: make([]StandardField, 0, 16),
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

	standardsFieldConfig := cfg.Fields.StandardsMap
	for fName, fValue := range allFields {
		v := util.AnyValueFromRaw(lineNo, fValue, cfg.Replace)

		fConfig, contains := standardsFieldConfig[fName]
		if contains {
			f := &StandardFieldT{Value: v, Config: fConfig}
			r.StandardFields = append(r.StandardFields, f)
		} else {
			r.OtherFields[fName] = v
		}
	}

	return r
}
