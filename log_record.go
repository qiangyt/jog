package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/gookit/goutil/strutil"
	"github.com/pkg/errors"
	"github.com/qiangyt/jog/config"
	"github.com/qiangyt/jog/util"
)

// StandardFieldT ...
type StandardFieldT struct {
	Value  string
	Config config.Field
}

// StandardField ...
type StandardField = *StandardFieldT

// LogRecordT ...
type LogRecordT struct {
	LineNo int

	Prefix         string
	StandardFields []StandardField
	OtherFields    map[string]interface{}
	Raw            string

	Unknown     bool
	StartupLine bool
}

// LogRecord .
type LogRecord = *LogRecordT

// PrintElement ...
func (i LogRecord) PrintElement(config Config, element config.Printable, builder *strings.Builder, a string) {
	if !element.IsEnabled() {
		return
	}

	var color util.Color
	if i.StartupLine {
		color = config.StartupLine.Color
	} else {
		color = element.GetColor(a)
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
		i.PrintElement(cfg, v, builder, strutil.MustString(fValue))
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
		i.PrintElement(cfg, f.Config, builder, f.Value)

		result[f.Config.Name] = builder.String()
	}
}

// AsFlatLine ...
func (i LogRecord) AsFlatLine(cfg Config) string {
	builder := &strings.Builder{}

	printStartLine := i.StartupLine && cfg.StartupLine.IsEnabled()
	startupLineColor := cfg.StartupLine.GetColor("")
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
		startupLineColor.Sprint(builder.String)
	}

	return builder.String()
}

func isStartupLine(cfg Config, raw string) bool {
	contains := cfg.StartupLine.Contains
	return len(contains) > 0 && strings.Contains(raw, contains)
}

// ParseRawLine ...
func ParseRawLine(cfg Config, lineNo int, rawLine string) LogRecord {
	r := &LogRecordT{
		LineNo:         lineNo,
		OtherFields:    make(map[string]interface{}),
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

	if lineNo == 360 {
		lineNo = 360
	}

	allFields := make(map[string]interface{})
	if err := json.Unmarshal([]byte(line), &allFields); err != nil {
		line = strings.ReplaceAll(line, "\\\"", "\"")
		if err := json.Unmarshal([]byte(line), &allFields); err != nil {
			log.Printf("failed to parse line %d: <%s>\n\treason %v\n", lineNo, rawLine, errors.Wrap(err, ""))
			return r
		}
	}

	r.Unknown = false

	standardsFieldConfig := cfg.Fields.StandardsMap
	for fName, fValue := range allFields {
		var v string

		alreadyNormalized := false

		if fValue != nil {
			kind := reflect.TypeOf(fValue).Kind()
			if kind == reflect.Map {
				json, err := json.MarshalIndent(fValue, "", "  ")
				if err != nil {
					log.Printf("line %v: failed to json format: %v\n", lineNo, fValue)
				} else {
					v = string(json)
				}
				alreadyNormalized = true
			} else {
				v = strutil.MustString(fValue)
			}

			if len(v) >= 1 {
				if v[:1] == "\"" || v[:1] == "'" {
					v = v[1:]
				}
			}
			if len(v) >= 1 {
				if v[len(v)-1:] == "\"" || v[len(v)-1:] == "'" {
					v = v[:len(v)-1]
				}
			}
			v = strutil.Replaces(v, cfg.Replace)

			if alreadyNormalized == false {
				var obj interface{}
				if err := json.Unmarshal([]byte(v), &obj); err == nil {
					json, err := json.MarshalIndent(obj, "", "  ")
					if err == nil {
						v = string(json)
					}
				}
			}
		}

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

// ProcessRawLine ...
func ProcessRawLine(cfg Config, lineNo int, rawLine string) {
	event := ParseRawLine(cfg, lineNo, rawLine)
	var line = event.AsFlatLine(cfg)
	if len(line) > 0 {
		fmt.Println(line)
	}
}

// ProcessLocalFile ...
func ProcessLocalFile(cfg Config, localFilePath string) {
	f, err := os.Open(localFilePath)
	if err != nil {
		panic(errors.Wrap(err, ""))
	}
	log.Printf("file is opened: %s\n", localFilePath)
	defer f.Close()

	ProcessReader(cfg, f)
}

// ProcessReader ...
func ProcessReader(cfg Config, reader io.Reader) {

	buf := bufio.NewReader(reader)

	for lineNo := 1; true; lineNo++ {
		rawLine, err := buf.ReadString('\n')
		len := len(rawLine)

		if len != 0 {
			// trim the tail \n
			if rawLine[len-1] == '\n' {
				rawLine = rawLine[:len-1]
			}
		}

		if err != nil {
			if err == io.EOF {
				log.Printf("got EOF, line %d\n", lineNo)
				ProcessRawLine(cfg, lineNo, rawLine)
				return
			}
			panic(errors.Wrapf(err, "failed to read line %d", lineNo))
		}

		ProcessRawLine(cfg, lineNo, rawLine)
	}
}
