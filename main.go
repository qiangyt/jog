package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gookit/color"
	"github.com/pkg/errors"
	"github.com/qiangyt/jog/jsonpath"
	"github.com/qiangyt/jog/util"
	"gopkg.in/yaml.v2"
)

// ParseConfigExpression ...
func ParseConfigExpression(expr string) (string, string, error) {
	arr := strings.Split(expr, "=")
	if len(arr) != 2 {
		return "", "", fmt.Errorf("invalid config item expression: <%s>", expr)
	}
	return arr[0], arr[1], nil
}

// ReadConfig ...
func ReadConfig(configFilePath string) Config {
	if len(configFilePath) == 0 {
		return ConfigWithDefaultYamlFile()
	}
	return ConfigWithYamlFile(configFilePath)
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
func SetConfigItem(cfg Config, m map[string]interface{}, configItemPath string, configItemValue string) {
	if err := jsonpath.Set(m, configItemPath, configItemValue); err != nil {
		panic(errors.Wrap(err, ""))
	}
	if err := cfg.FromMap(m); err != nil {
		panic(errors.Wrap(err, ""))
	}
}

func main() {
	ok, cmdLine := ParseCommandLine()
	if !ok {
		return
	}

	if !cmdLine.Debug {
		defer func() {
			if p := recover(); p != nil {
				color.Red.Printf("%v\n\n", p)
				os.Exit(1)
				return
			}
		}()
	}

	logFile := util.InitLogger()
	defer logFile.Close()

	cfg := ReadConfig(cmdLine.ConfigFilePath)

	if len(cmdLine.ConfigItemPath) > 0 {
		m := cfg.ToMap()
		if len(cmdLine.ConfigItemValue) > 0 {
			SetConfigItem(cfg, m, cmdLine.ConfigItemPath, cmdLine.ConfigItemValue)
		} else {
			PrintConfigItem(m, cmdLine.ConfigItemPath)
			return
		}
	}

	if len(cmdLine.LogFilePath) == 0 {
		log.Println("read JSON log lines from stdin")
		ProcessReader(cfg, os.Stdin, 1)
	} else {
		log.Printf("processing local JSON log file: %s\n", cmdLine.LogFilePath)
		ProcessLocalFile(cfg, cmdLine.FollowMode, cmdLine.LogFilePath)
	}

	fmt.Println()
}
