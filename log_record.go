package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/gookit/goutil/strutil"
	"github.com/pkg/errors"
	"github.com/qiangyt/jog/config"
	"github.com/qiangyt/jog/util"
)

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

	element.PrintTo(color, builder, a)
}

// PopulateOtherFields ...
func (i LogRecord) PopulateOtherFields(cfg config.Configuration, unknownFields map[string]util.AnyValue, implicitStandardFields map[string]FieldValue, result map[string]string) {
	if !cfg.HasOthersFieldInPattern {
		return
	}

	nameElement := cfg.Fields.Others.Name
	separatorElement := cfg.Fields.Others.Separator
	unknownFieldValueElement := cfg.Fields.Others.Value

	// sort field names
	var fNames []string
	for fName := range unknownFields {
		fNames = append(fNames, fName)
	}
	for fName := range implicitStandardFields {
		fNames = append(fNames, fName)
	}
	sort.Strings(fNames)

	builder := &strings.Builder{}
	first := true

	for _, fName := range fNames {
		fValueUnknown := unknownFields[fName]
		if fValueUnknown != nil {
			if !first {
				builder.WriteString(", ")
			}
			first = false

			i.PrintElement(cfg, nameElement, builder, fName)
			i.PrintElement(cfg, separatorElement, builder, "=")
			i.PrintElement(cfg, unknownFieldValueElement, builder, fValueUnknown.String())
		} else {
			fValueImplicit := implicitStandardFields[fName]

			if !fValueImplicit.Config.IsEnabled() {
				continue
			}

			if !first {
				builder.WriteString(", ")
			}
			first = false

			i.PrintElement(cfg, nameElement, builder, fName)
			i.PrintElement(cfg, separatorElement, builder, "=")
			i.PrintElement(cfg, fValueImplicit.Config, builder, fValueImplicit.Output)
		}
	}

	result["others"] = builder.String()
}

// PopulateExplicitStandardFields ...
func (i LogRecord) PopulateExplicitStandardFields(cfg config.Configuration, explicitStandardFields map[string]FieldValue, result map[string]string) {
	for _, f := range explicitStandardFields {
		builder := &strings.Builder{}
		i.PrintElement(cfg, f.Config, builder, f.Output)

		result[f.Config.Name] = builder.String()
	}
}

// AsFlatLine ...
func (i LogRecord) AsFlatLine(cfg config.Configuration) string {
	builder := &strings.Builder{}

	i.PrintElement(cfg, cfg.LineNo, builder, fmt.Sprintf("%-6v ", i.LineNo))

	if i.Unknown {
		i.PrintElement(cfg, cfg.UnknownLine, builder, i.Raw)
	} else {
		if len(i.Prefix) > 0 {
			i.PrintElement(cfg, cfg.Prefix, builder, strutil.MustString(i.Prefix))
		}

		explicitStandardFields, implicitStandardFields := i.ExtractStandardFields(cfg)

		result := make(map[string]string)

		i.PopulateOtherFields(cfg, i.UnknownFields, implicitStandardFields, result)
		i.PopulateExplicitStandardFields(cfg, explicitStandardFields, result)

		builder.WriteString(os.Expand(cfg.Pattern, func(fieldName string) string {
			return result[fieldName]
		}))
	}

	if i.StartupLine {
		startupLine := &strings.Builder{}
		cfg.StartupLine.PrintTo(cfg.StartupLine.Color, startupLine, builder.String())
		builder = startupLine
	}

	return builder.String()
}

// ExtractStandardFields ...
func (i LogRecord) ExtractStandardFields(cfg config.Configuration) (map[string]FieldValue, map[string]FieldValue) {
	explicts := make(map[string]FieldValue)
	implicits := make(map[string]FieldValue)

	for n, f := range i.StandardFields {
		if cfg.HasFieldInPattern(n) {
			explicts[n] = f
		} else {
			implicits[n] = f
		}
	}

	return explicts, implicits
}

// MatchesLevelFilter ...
func (i LogRecord) MatchesLevelFilter(cfg config.Configuration, levelFilters []config.Enum) bool {
	levelFieldValue := i.StandardFields["level"]
	if levelFieldValue != nil {
		levelFieldEnum := levelFieldValue.enumValue
		for _, levelFilter := range levelFilters {
			if levelFieldEnum == levelFilter {
				return true
			}
		}
	}
	return false
}

// MatchesTimestampFilter ...
func (i LogRecord) MatchesTimestampFilter(cfg config.Configuration, beforeFilter *time.Time, afterFilter *time.Time) bool {
	timestampFieldValue := i.StandardFields["timestamp"]
	if timestampFieldValue != nil {
		timestampValue := timestampFieldValue.timeValue

		beforeMatches := true
		if beforeFilter != nil {
			beforeMatches = timestampValue.Before(*beforeFilter) || timestampValue.Equal(*beforeFilter)
		}

		afterMatches := true
		if afterFilter != nil {
			afterMatches = timestampValue.After(*afterFilter) || timestampValue.Equal(*afterFilter)
		}

		return beforeMatches && afterMatches
	}
	return false
}

// MatchesFilters ...
func (i LogRecord) MatchesFilters(cfg config.Configuration, options Options) bool {
	levelFilters := options.GetLevelFilters()

	if len(options.levelFilters) > 0 {
		if !i.MatchesLevelFilter(cfg, levelFilters) {
			return false
		}
	}

	if options.HasTimestampFilter() {
		if !i.MatchesTimestampFilter(cfg, options.BeforeFilter, options.AfterFilter) {
			return false
		}
	}

	return true
}

func isStartupLine(cfg config.Configuration, raw string) bool {
	contains := cfg.StartupLine.Contains
	return len(contains) > 0 && strings.Contains(raw, contains)
}

func tryToParseUsingGrok(cfg config.Configuration, options Options, lineNo int, line string) (matchesGrok bool, prefix string, standardFields map[string]FieldValue, unknownFields map[string]util.AnyValue) {
	prefix = ""
	standardFields = map[string]FieldValue{}
	unknownFields = map[string]util.AnyValue{}

	if options.isGrokEnabled() == false {
		matchesGrok = false
		return
	}

	standardsFieldConfig := cfg.Fields.StandardsWithAllAliases
	grok := cfg.Grok

	for _, pattern := range options.GrokPatterns {
		fields, err := grok.Parse(pattern, line)
		if err != nil {
			//TODO: debug log
		} else {
			for fName, fValue := range fields {
				v := util.AnyValueFromRaw(lineNo, fValue, cfg.Replace)

				fConfig, contains := standardsFieldConfig[fName]
				if contains {
					fName = fConfig.Name // normalize field name
					standardFields[fName] = NewFieldValue(cfg, options, fConfig, v)
				} else {
					unknownFields[fName] = v
				}
			}

			if len(standardFields) == 0 {
				continue
			}

			matchesGrok = true

			for _, matchesFieldName := range grok.MatchesFields {
				if _, contains := standardFields[matchesFieldName]; contains == false {
					matchesGrok = false
					break
				}
			}

			if matchesGrok {
				return
			}
		}
	}

	return
}

func tryToParseAsJSON(cfg config.Configuration, options Options, lineNo int, line string) (isJSON bool, prefix string, standardFields map[string]FieldValue, unknownFields map[string]util.AnyValue) {
	prefix = ""
	standardFields = map[string]FieldValue{}
	unknownFields = map[string]util.AnyValue{}

	posOfLeftBracket := strings.IndexByte(line, '{')
	if posOfLeftBracket < 0 {
		log.Printf("line %d is not JSON line: <%s>\n", lineNo, line)
		isJSON = false
		return
	}

	if posOfLeftBracket > 0 {
		prefix = line[:posOfLeftBracket]
		line = line[posOfLeftBracket:]
	}

	allFields := make(map[string]interface{})

	if err := json.Unmarshal([]byte(line), &allFields); err != nil {
		log.Printf("parse round 1 failed: line %d: <%s>\n\treason %v\n", lineNo, line, errors.Wrap(err, ""))
		line = strings.ReplaceAll(line, "\\\"", "\"")
		line = strings.ReplaceAll(line, "\t", "    ")
		if err := json.Unmarshal([]byte(line), &allFields); err != nil {
			log.Printf("parse round 2 failed: line %d: <%s>\n\treason %v\n", lineNo, line, errors.Wrap(err, ""))
			isJSON = false
			return
		}
	}

	standardsFieldConfig := cfg.Fields.StandardsWithAllAliases
	for fName, fValue := range allFields {
		v := util.AnyValueFromRaw(lineNo, fValue, cfg.Replace)

		fConfig, contains := standardsFieldConfig[fName]
		if contains {
			fName = fConfig.Name // normalize field name
			standardFields[fName] = NewFieldValue(cfg, options, fConfig, v)
		} else {
			unknownFields[fName] = v
		}
	}

	isJSON = true
	return
}

// ParseAsRecord ...
func ParseAsRecord(cfg config.Configuration, options Options, lineNo int, rawLine string) LogRecord {
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

	var isJSON bool
	var matchesGrok bool

	isJSON, r.Prefix, r.StandardFields, r.UnknownFields = tryToParseAsJSON(cfg, options, lineNo, line)
	if isJSON {
		r.Unknown = false
	} else {
		matchesGrok, r.Prefix, r.StandardFields, r.UnknownFields = tryToParseUsingGrok(cfg, options, lineNo, line)
		if matchesGrok {
			r.Unknown = false
		} else {
			r.Unknown = true
		}
	}

	return r
}
