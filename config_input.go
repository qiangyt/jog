package main

import (
	"strings"

	"github.com/gookit/goutil/strutil"
	"gopkg.in/yaml.v2"
)

// FieldNameT ...
type FieldNameT struct {
	yaml.Unmarshaler
	yaml.Marshaler

	Names []string
}

// FieldName ...
type FieldName = *FieldNameT

// NewFieldName ...
func NewFieldName(namesText string) FieldName {
	return &FieldNameT{
		Names: strutil.Split(namesText, ","),
	}
}

// UnmarshalYAML ...
func (me FieldName) UnmarshalYAML(unmarshal func(interface{}) error) error {
	namesText := ""
	err := unmarshal(&namesText)
	if err != nil {
		return err
	}

	me.Names = strutil.Split(namesText, ",")

	return nil
}

// MarshalYAML ...
func (me FieldName) MarshalYAML() (interface{}, error) {
	return strings.Join(me.Names, ", "), nil
}

// FieldNamesConfigT ...
type FieldNamesConfigT struct {
	App        FieldName
	Timestamp  FieldName
	Version    FieldName
	Message    FieldName
	Logger     FieldName
	Thread     FieldName
	Level      FieldName
	StackTrace FieldName `yaml:"stack-trace"`
	PID        FieldName `yaml:"pid"`
	Host       FieldName
	File       FieldName
	Method     FieldName
	Line       FieldName
}

// FieldNamesConfig ...
type FieldNamesConfig = *FieldNamesConfigT

// DefaultFieldNamesConfig ...
func DefaultFieldNamesConfig() FieldNamesConfig {
	return &FieldNamesConfigT{
		App:        NewFieldName("name, Name, app, App, @name, @Name, @app, @App"),
		Timestamp:  NewFieldName("time, Time, timestamp, Timestamp, @time, @Time, @timestamp, @Timestamp"),
		Version:    NewFieldName("version, Version, @version, @Version"),
		Message:    NewFieldName("msg, message, Message, @msg, @message, @Message"),
		Logger:     NewFieldName("id, Id, ID, logger_name, logger-name, loggerName, LoggerName, logger, Logger, @id, @Id, @ID, @logger_name, @logger-name, @loggerName, @LoggerName, @logger, @Logger"),
		Thread:     NewFieldName("thread_name, thread-name, threadName, ThreadName, thread, Thread, @thread, @Thread"),
		Level:      NewFieldName("level, Level, @level, @Level"),
		StackTrace: NewFieldName("stack_trace, stack-trace, stackTrace, StackTrace, stack, Stack, @stack_trace, @stack-trace, @stackTrace, @StackTrace, @stack, @Stack"),
		PID:        NewFieldName("pid, PID, @pid, @PID"),
		Host:       NewFieldName("host, Host, @host, @Host, hostname, Hostname, hostName, HostName, @Hostname, @Hostname, @hostName, @HostName"),
		File:       NewFieldName("file, File, @file, @File"),
		Method:     NewFieldName("method, Method, @method, @Method"),
		Line:       NewFieldName("line, Line, @line, @Line"),
	}
}

// InputConfigT ...
type InputConfigT struct {
	FieldNames            FieldNamesConfig `yaml:"field-names"`
	IgnoreConversionError bool             `yaml:"ignore-conversion-error"`
}

// InputConfig ...
type InputConfig = *InputConfigT

// DefaultInputConfig ...
func DefaultInputConfig() InputConfig {
	return &InputConfigT{
		FieldNames:            DefaultFieldNamesConfig(),
		IgnoreConversionError: true,
	}
}
