package convert

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"github.com/qiangyt/jog/jsonpath"
	"github.com/qiangyt/jog/util"
	"gopkg.in/yaml.v2"
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

func (i ConvertContext) LoadConfig() bool {
	var cfg Config
	var options = i.Options()

	if len(options.ConfigFilePath) == 0 {
		cfg = NewConfigWithDefaultYamlFile(i)
	} else {
		cfg = NewConfigWithYamlFile(i, options.ConfigFilePath)
	}
	i.config = cfg

	if len(options.ConfigItemPath) > 0 {
		m := cfg.ToMap()
		if len(options.ConfigItemValue) > 0 {
			i.SetConfigItem(m, options.ConfigItemPath, options.ConfigItemValue)
		} else {
			i.PrintConfigItem(m, options.ConfigItemPath)
			return false
		}
	}

	if cfg.LevelField != nil {
		options.InitLevelFilters(cfg.LevelField.Enums)
	}
	if cfg.TimestampField != nil {
		options.InitTimestampFilters(i)
	}

	options.InitGroks(cfg)

	return true
}

func (i ConvertContext) PrintConfigItem(m map[string]interface{}, configItemPath string) {
	item, err := jsonpath.Get(m, configItemPath)
	if err != nil {
		panic(errors.Wrap(err, ""))
	}
	out, err := yaml.Marshal(item)
	if err != nil {
		panic(errors.Wrap(err, ""))
	}
	fmt.Print(string(out))
}

func (i ConvertContext) SetConfigItem(m map[string]interface{}, configItemPath string, configItemValue string) {
	if err := jsonpath.Set(m, configItemPath, configItemValue); err != nil {
		panic(errors.Wrap(err, ""))
	}
	if err := i.Config().FromMap(m); err != nil {
		panic(errors.Wrap(err, ""))
	}
}

func (i ConvertContext) Config() Config {
	return i.config
}
