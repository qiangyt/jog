package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/qiangyt/jog/config"
	"github.com/qiangyt/jog/util"
	"gopkg.in/yaml.v2"
)

// ConfigT ...
type ConfigT struct {
	// TODO: configurable
	Replacer    map[string]string
	Pattern     string
	StartupLine config.StartupLine `yaml:"startup-line"`
	LineNo      config.Element     `yaml:"line-no"`
	UnknownLine config.Element     `yaml:"unknown-line"`
	Prefix      config.Prefix
	Fields      config.FieldMap
}

// Config ...
type Config = *ConfigT

func lookForConfigFile(dir string) string {
	log.Printf("looking for config files in: %s\n", dir)
	r := filepath.Join(dir, ".jog.yaml")
	if util.FileExists(r) {
		return r
	}
	r = filepath.Join(dir, ".jog.yml")
	if util.FileExists(r) {
		return r
	}
	return ""
}

// DetermineConfigFilePath return (file path)
func DetermineConfigFilePath() string {
	dir := util.ExeDirectory()
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
		return ConfigWithYaml(config.DefaultYAML)
	}

	log.Printf("Config file: %s\n", path)
	return ConfigWithYaml(path)
}

// ConfigWithYamlFile ...
func ConfigWithYamlFile(path string) Config {
	log.Printf("Config file: %s\n", path)

	yamlText := string(util.ReadFile(path))
	return ConfigWithYaml(yamlText)
}

// ConfigWithYaml ...
func ConfigWithYaml(yamlText string) Config {
	r := &ConfigT{
		Replacer: map[string]string{
			"\\\"": "\"",
			"\\'":  "'",
			"\\\n": "\n",
			"\\\r": "",
			"\\\t": "\t",
		},
		Pattern:     "",
		StartupLine: &config.StartupLineT{},
		LineNo:      &config.ElementT{},
		UnknownLine: &config.ElementT{},
		Prefix:      &config.PrefixT{},
		Fields:      &config.FieldMapT{},
	}
	if err := yaml.Unmarshal([]byte(yamlText), &r); err != nil {
		panic(errors.Wrap(err, "failed to unmarshal yaml: \n"+yamlText))
	}
	return r
}
