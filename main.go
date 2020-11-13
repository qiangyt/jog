package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gookit/color"
	"github.com/pkg/errors"
	"github.com/qiangyt/jog/config"
	"github.com/qiangyt/jog/jsonpath"
	"github.com/qiangyt/jog/util"
	"gopkg.in/yaml.v2"
)

//go:generate go run script/include_static.go

// ParseConfigExpression ...
func ParseConfigExpression(expr string) (string, string, error) {
	arr := strings.Split(expr, "=")
	if len(arr) != 2 {
		return "", "", fmt.Errorf("invalid config item expression: <%s>", expr)
	}
	return arr[0], arr[1], nil
}

// ReadConfig ...
func ReadConfig(configFilePath string) config.Configuration {
	if len(configFilePath) == 0 {
		return config.WithDefaultYamlFile()
	}
	return config.WithYamlFile(configFilePath)
}

// PrintConfigItem ...
func PrintConfigItem(m map[string]interface{}, configItemPath string) {
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

// SetConfigItem ...
func SetConfigItem(cfg config.Configuration, m map[string]interface{}, configItemPath string, configItemValue string) {
	if err := jsonpath.Set(m, configItemPath, configItemValue); err != nil {
		panic(errors.Wrap(err, ""))
	}
	if err := cfg.FromMap(m); err != nil {
		panic(errors.Wrap(err, ""))
	}
}

func main() {
	config.InitDefaultGrokLibraryDir()

	ok, options := OptionsWithCommandLine()
	if !ok {
		return
	}

	if !options.Debug {
		defer func() {
			if p := recover(); p != nil {
				color.Red.Printf("%v\n\n", p)
				os.Exit(1)
				return
			}
		}()
	}

	logFile := util.InitLogger(config.JogHomeDir(true))
	defer logFile.Close()

	cfg := ReadConfig(options.ConfigFilePath)

	if len(options.ConfigItemPath) > 0 {
		m := cfg.ToMap()
		if len(options.ConfigItemValue) > 0 {
			SetConfigItem(cfg, m, options.ConfigItemPath, options.ConfigItemValue)
		} else {
			PrintConfigItem(m, options.ConfigItemPath)
			return
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

	fmt.Println()
}
