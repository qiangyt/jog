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

func readConfig(ctx util.JogContext, configFilePath string) Config {
	if len(configFilePath) == 0 {
		return NewConfigWithDefaultYamlFile(ctx)
	}
	return NewConfigWithYamlFile(ctx, configFilePath)
}

func printConfigItem(m map[string]interface{}, configItemPath string) {
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

func setConfigItem(cfg Config, m map[string]interface{}, configItemPath string, configItemValue string) {
	if err := jsonpath.Set(m, configItemPath, configItemValue); err != nil {
		panic(errors.Wrap(err, ""))
	}
	if err := cfg.FromMap(m); err != nil {
		panic(errors.Wrap(err, ""))
	}
}

// NewConvertContext ...
func NewConvertContext(jogHomeDir string, version string) (util.LogFile, util.JogContext) {
	lf := util.NewLogFile(filepath.Join(jogHomeDir, "convert.log"))

	logger := log.With(log.NewStdLogger(lf.File()),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
	)

	ctx := util.NewJogContext(context.TODO())
	ctx.WithLogger(logger)
	ctx.LogInfo("------------------------------------------------", "version", version, "pid", os.Getpid())

	return lf, ctx
}

func Main(args []string, version string) Options {
	util.InitDefaultGrokLibraryDir()

	ok, options := NewOptionsWithCommandLine(args)
	if !ok {
		return nil
	}

	logFile, ctx := NewConvertContext(util.JogHomeDir(true), version)
	defer logFile.Close()

	cfg := readConfig(ctx, options.ConfigFilePath)

	if len(options.ConfigItemPath) > 0 {
		m := cfg.ToMap()
		if len(options.ConfigItemValue) > 0 {
			setConfigItem(cfg, m, options.ConfigItemPath, options.ConfigItemValue)
		} else {
			printConfigItem(m, options.ConfigItemPath)
			return options
		}
	}

	if cfg.LevelField != nil {
		options.InitLevelFilters(cfg.LevelField.Enums)
	}
	if cfg.TimestampField != nil {
		options.InitTimestampFilters(ctx, cfg.TimestampField)
	}

	options.InitGroks(cfg)

	if len(options.LogFilePath) == 0 {
		ctx.LogInfo("read JSON log lines from stdin")
		ProcessReader(ctx, cfg, options, os.Stdin, 1)
	} else {
		ctx.LogInfo("processing local JSON log file", "logFilePath", options.LogFilePath)
		ProcessLocalFile(ctx, cfg, options, options.FollowMode, options.LogFilePath)
	}

	return options
}
