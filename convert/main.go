package convert

import (
	"fmt"
	"log"
	"os"

	"github.com/pkg/errors"
	"github.com/qiangyt/jog/jsonpath"
	"github.com/qiangyt/jog/util"
	"gopkg.in/yaml.v2"
)

func readConfig(configFilePath string) Config {
	if len(configFilePath) == 0 {
		return NewConfigWithDefaultYamlFile()
	}
	return NewConfigWithYamlFile(configFilePath)
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

func Main(args []string) Options {
	util.InitDefaultGrokLibraryDir()

	ok, options := NewOptionsWithCommandLine(args)
	if !ok {
		return nil
	}

	logFile := util.InitLogLogger(util.JogHomeDir(true))
	defer logFile.Close()

	cfg := readConfig(options.ConfigFilePath)

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
		options.InitTimestampFilters(cfg.TimestampField)
	}

	options.InitGroks(cfg)

	if len(options.LogFilePath) == 0 {
		log.Println("read JSON log lines from stdin")
		ProcessReader(cfg, options, os.Stdin, 1)
	} else {
		log.Printf("processing local JSON log file: %s\n", options.LogFilePath)
		ProcessLocalFile(cfg, options, options.FollowMode, options.LogFilePath)
	}

	return options
}
