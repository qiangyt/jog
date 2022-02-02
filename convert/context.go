package convert

import (
	"context"
	"os"
	"path/filepath"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/qiangyt/jog/util"
)

type ConvertContextT struct {
	util.JogContextT
	config  Config
	options Options
	logFile util.LogFile
}

type ConvertContext = *ConvertContextT

// NewConvertContext ...
func NewConvertContext(options Options, jogHomeDir string, version string) ConvertContext {
	logFile := util.NewLogFile(filepath.Join(jogHomeDir, "convert.log"))

	logger := log.With(log.NewStdLogger(logFile.File()),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
	)

	ctx := util.NewJogContext(context.TODO())
	ctx.WithLogger(logger)
	ctx.LogInfo("------------------------------------------------", "version", version, "pid", os.Getpid())

	return &ConvertContextT{
		JogContextT: ctx,
		logFile:     logFile,
		options:     options,
	}
}

func (i ConvertContext) Close() {
	i.logFile.Close()
}

func (i ConvertContext) Options() Options {
	return i.options
}

func (i ConvertContext) LoadConfig() Config {
	configFilePath := i.options.ConfigFilePath

	if len(configFilePath) == 0 {
		i.config = NewConfigWithDefaultYamlFile(i)
	} else {
		i.config = NewConfigWithYamlFile(i, configFilePath)
	}

	return i.config
}

func (i ConvertContext) Config() Config {
	return i.config
}
