package main

import (
	"github.com/gookit/goutil/strutil"
	"github.com/pkg/errors"
)

// GenerialMediatorT implements LogMediator interface
type GenerialMediatorT struct {
}

// GenerialMediator is pointer of GenerialMediatorT
type GenerialMediator = *GenerialMediatorT

func extractField(cfg Config, fields map[string]interface{}, fieldName string, amountOfFieldsPopulated *int) interface{} {
	r, has := fields[fieldName]
	if has {
		delete(fields, fieldName)
		(*amountOfFieldsPopulated)++
		return r
	}
	return nil
}

func extractFieldString(cfg Config, fields map[string]interface{}, fieldName string, amountOfFieldsPopulated *int) string {
	i := extractField(cfg, fields, fieldName, amountOfFieldsPopulated)
	if i == nil {
		return ""
	}

	r, err := strutil.ToString(i)
	if !cfg.Input.IgnoreConversionError {
		panic(errors.Wrapf(err, "failed to convert '%v' to string", i))
	}
	return r
}

// PopulateFields populates field into the log event
func (me GenerialMediator) PopulateFields(cfg Config, event LogEvent, fields map[string]interface{}) int {
	r := 0
	fieldNamesConfig := cfg.Input.FieldNames

	event.Timestamp = extractFieldString(cfg, fields, fieldNamesConfig.Timestamp, &r)
	event.Version = extractFieldString(cfg, fields, fieldNamesConfig.Version, &r)
	event.Message = extractFieldString(cfg, fields, fieldNamesConfig.Message, &r)
	event.Logger = extractFieldString(cfg, fields, fieldNamesConfig.Logger, &r)
	event.Thread = extractFieldString(cfg, fields, fieldNamesConfig.Thread, &r)
	event.Level = extractFieldString(cfg, fields, fieldNamesConfig.Level, &r)
	event.StackTrace = extractFieldString(cfg, fields, fieldNamesConfig.StackTrace, &r)
	event.PID = extractFieldString(cfg, fields, fieldNamesConfig.PID, &r)
	event.Host = extractFieldString(cfg, fields, fieldNamesConfig.Host, &r)
	event.File = extractFieldString(cfg, fields, fieldNamesConfig.File, &r)
	event.Method = extractFieldString(cfg, fields, fieldNamesConfig.Method, &r)
	event.Line = extractFieldString(cfg, fields, fieldNamesConfig.Line, &r)

	event.Others = fields

	return r
}
