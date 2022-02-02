package convert

import (
	"fmt"

	"os"

	"github.com/pkg/errors"
	"github.com/qiangyt/jog/jsonpath"
	"github.com/qiangyt/jog/util"
	"gopkg.in/yaml.v2"
)

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

func Main(args []string, version string) Options {
	util.InitDefaultGrokLibraryDir()

	ok, options := NewOptionsWithCommandLine(args)
	if !ok {
		return nil
	}

	ctx := NewConvertContext(options, util.JogHomeDir(true), version)
	defer ctx.Close()

	cfg := ctx.LoadConfig()

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
		options.InitTimestampFilters(ctx)
	}

	options.InitGroks(cfg)

	if len(options.LogFilePath) == 0 {
		ctx.LogInfo("read JSON log lines from stdin")
		ProcessReader(ctx, os.Stdin, 1)
	} else {
		ctx.LogInfo("processing local JSON log file", "logFilePath", options.LogFilePath)
		ProcessLocalFile(ctx)
	}

	return options
}
