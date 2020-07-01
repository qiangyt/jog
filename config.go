package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// OutputConfigT ...
type OutputConfigT struct {
	Pattern            string
	CompressLoggerName bool `yaml:"compress-logger-name"`
	Colors             OutputColorsConfig
	StartedLine        string `yaml:"started-line"`
}

// OutputConfig ...
type OutputConfig = *OutputConfigT

// ConfigT ...
type ConfigT struct {
	Output OutputConfigT
}

// Config ...
type Config = *ConfigT

func lookForConfigFile(dir string) string {
	log.Printf("looking for config files in: %s\n", dir)
	r := filepath.Join(dir, ".j2log.yaml")
	if FileExists(r) {
		return r
	}
	r = filepath.Join(dir, ".j2log.yml")
	if FileExists(r) {
		return r
	}
	return ""
}

// DetermineConfigFilePath return (file path)
func DetermineConfigFilePath() string {
	for i, arg := range os.Args {
		if arg == "-c" {
			log.Println("got -c option")
			if i == len(os.Args)-1 {
				panic(fmt.Errorf("missing config file for -c option"))
			}
			return os.Args[i+1]
		}
	}

	dir := ExeDirectory()
	r := lookForConfigFile(dir)
	if len(r) != 0 {
		return r
	}

	dir = os.ExpandEnv("${HOME}")
	return lookForConfigFile(dir)
}

// ConfigWithDefaultYamlFile ...
func ConfigWithDefaultYamlFile() Config {
	path := DetermineConfigFilePath()

	if len(path) == 0 {
		log.Println("Config file not found, take default config")
		return ConfigWithYaml(ConfigDefaultYAML)
	}

	log.Printf("Config file: %s\n", path)
	return ConfigWithYaml(path)
}

// ConfigWithYamlFile ...
func ConfigWithYamlFile(path string) Config {
	log.Printf("Config file: %s\n", path)

	yamlText := string(ReadFile(path))
	return ConfigWithYaml(yamlText)
}

// ConfigWithYaml ...
func ConfigWithYaml(yamlText string) Config {
	var r ConfigT
	if err := yaml.UnmarshalStrict([]byte(yamlText), &r); err != nil {
		panic(errors.Wrap(err, "failed to unmarshal yaml: \n"+yamlText))
	}
	return &r
}
