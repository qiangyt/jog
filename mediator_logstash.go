package main

import "strings"

// LogstashMediatorT implements LogMediator interface
type LogstashMediatorT struct {
}

// LogstashMediator is pointer of LogstashMediatorT
type LogstashMediator = *LogstashMediatorT

// PopulateFields populates field into the log event
func (me LogstashMediator) PopulateFields(cfg Config, event LogEvent) int {
	amountOfFieldsPopulated := 0

	for fieldName, fieldValue := range event.All {
		if fieldName == "@timestamp" {
			event.Timestamp = fieldValue.(string)
			amountOfFieldsPopulated++
			continue
		}
		if fieldName == "@version" {
			event.Version = fieldValue.(string)
			amountOfFieldsPopulated++
			continue
		}
		if fieldName == "message" {
			event.Message = fieldValue.(string)
			amountOfFieldsPopulated++
			continue
		}
		if fieldName == "logger_name" {
			event.Logger = fieldValue.(string)
			amountOfFieldsPopulated++
			continue
		}
		if fieldName == "thread_name" {
			event.Thread = fieldValue.(string)
			amountOfFieldsPopulated++
			continue
		}
		if fieldName == "level" {
			event.Level = strings.ToLower(fieldValue.(string))
			amountOfFieldsPopulated++
			continue
		}
		if fieldName == "level_value" {
			amountOfFieldsPopulated++
			// skip
			continue
		}
		if fieldName == "stack_trace" {
			event.StackTrace = fieldValue.(string)
			amountOfFieldsPopulated++
			continue
		}

		event.Others[fieldName] = fieldValue
	}

	return amountOfFieldsPopulated
}
