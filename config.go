package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// ConfigT ...
type ConfigT struct {
	Output OutputConfig
	Input  InputConfig
}

// DefaultConfig ...
func DefaultConfig() Config {
	return &ConfigT{
		Output: DefaultOutputConfig(),
		Input:  DefaultInputConfig(),
	}
}

// Config ...
type Config = *ConfigT

func lookForConfigFile(dir string) string {
	log.Printf("looking for config files in: %s\n", dir)
	r := filepath.Join(dir, ".jog.yaml")
	if FileExists(r) {
		return r
	}
	r = filepath.Join(dir, ".jog.yml")
	if FileExists(r) {
		return r
	}
	return ""
}

// DetermineConfigFilePath return (file path)
func DetermineConfigFilePath() string {
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
	r := DefaultConfig()
	if err := yaml.UnmarshalStrict([]byte(yamlText), &r); err != nil {
		panic(errors.Wrap(err, "failed to unmarshal yaml: \n"+yamlText))
	}
	return r
}
