package _util

import (
	"context"
	"os"

	"github.com/go-kratos/kratos/v2/log"
)

type JogContextT struct {
	context.Context

	logger log.Logger
}

type JogContext = *JogContextT

func NewJogContext(parent context.Context) JogContextT {
	if parent == nil {
		parent = context.TODO()
	}
	return JogContextT{Context: parent}
}

// NewConvertContext ...
func (i JogContext) WithLogger(logger log.Logger, version string) {
	i.logger = logger
	i.LogInfo("------------------------------------------------", "version", version, "pid", os.Getpid())
}

func NewTestContext() JogContext {
	r := NewJogContext(nil)
	r.WithLogger(log.With(log.DefaultLogger), "test")
	return &r
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
