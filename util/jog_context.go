package util

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

type JogContextT struct {
	context.Context

	logger log.Logger
}

type JogContext = *JogContextT

func NewJogContext(parent context.Context) JogContext {
	if parent == nil {
		parent = context.TODO()
	}
	return &JogContextT{Context: parent}
}

// NewConvertContext ...
func (i JogContext) WithLogger(logger log.Logger) {
	i.logger = logger
}

func NewTestContext() JogContext {
	r := NewJogContext(nil)
	r.WithLogger(log.With(log.DefaultLogger))
	return r
}

func (i JogContext) LogInfo(msg string, keyvals ...interface{}) error {
	keyvals = append(keyvals, "msg", msg)
	return i.logger.Log(log.LevelInfo, keyvals...)
}

func (i JogContext) LogDebug(msg string, keyvals ...interface{}) error {
	keyvals = append(keyvals, "msg", msg)
	return i.logger.Log(log.LevelDebug, keyvals...)
}

func (i JogContext) LogWarn(msg string, keyvals ...interface{}) error {
	keyvals = append(keyvals, "msg", msg)
	return i.logger.Log(log.LevelWarn, keyvals...)
}

func (i JogContext) LogError(msg string, keyvals ...interface{}) error {
	keyvals = append(keyvals, "msg", msg)
	return i.logger.Log(log.LevelError, keyvals...)
}

func (i JogContext) LogFatal(msg string, keyvals ...interface{}) error {
	keyvals = append(keyvals, "msg", msg)
	return i.logger.Log(log.LevelFatal, keyvals...)
}
