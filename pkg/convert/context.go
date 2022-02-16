package convert

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"github.com/qiangyt/jog/pkg/jsonpath"
	_log "github.com/qiangyt/jog/pkg/log"
	_util "github.com/qiangyt/jog/pkg/util"
	"gopkg.in/yaml.v2"
)

type ConvertContextT struct {
	_util.JogContextT
	config  Config
	options Options
	logFile _log.File
}

type ConvertContext = *ConvertContextT

// NewConvertContext ...
func NewConvertContext(options Options, jogHomeDir string, version string) ConvertContext {
	logFile := _log.NewFile(filepath.Join(jogHomeDir, "convert.log"))

	logger := log.With(log.NewStdLogger(logFile.File()),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
	)

	ctx := _util.NewJogContext(context.TODO())
	ctx.WithLogger(logger, version)

	return &ConvertContextT{
		JogContextT: ctx,
		logFile:     logFile,
		options:     options,
	}
}

func (i ConvertContext) Close() {
	if i.logFile != nil {
		i.logFile.Close()
	}
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
