package main

// LogstashFormatT implements LogFormat interface
type LogstashFormatT struct {
}

// LogstashFormat is pointer of LogstashFormatT
type LogstashFormat = *LogstashFormatT

// Parse the log event
func (me LogstashFormat) Parse(event LogEvent) error {
	for fieldName, fieldValue := range event.All {
		if fieldName == "@timestamp" {
			event.Timestamp = fieldValue.(string)
			continue
		}
		if fieldName == "@version" {
			event.Version = fieldValue.(string)
			continue
		}
		if fieldName == "message" {
			event.Message = fieldValue.(string)
			continue
		}
		if fieldName == "logger_name" {
			event.Logger = fieldValue.(string)
			continue
		}
		if fieldName == "thread_name" {
			event.Thread = fieldValue.(string)
			continue
		}
		if fieldName == "level" {
			event.Level = fieldValue.(string)
			continue
		}
		if fieldName == "level_value" {
			// skip
			continue
		}
		if fieldName == "stack_trace" {
			event.StackTrace = fieldValue.(string)
			continue
		}

		event.Others[fieldName] = fieldValue
	}

	return nil
}
