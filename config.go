package main

import (
	"log"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/qiangyt/jog/config"
	"github.com/qiangyt/jog/util"
	"gopkg.in/yaml.v2"
)

// ConfigT ...
type ConfigT struct {
	// TODO: configurable
	Colorization bool
	Replace      map[string]string
	Pattern      string
	StartupLine  config.StartupLine `yaml:"startup-line"`
	LineNo       config.Element     `yaml:"line-no"`
	UnknownLine  config.Element     `yaml:"unknown-line"`
	Prefix       config.Prefix
	Fields       config.FieldMap
}

// Config ...
type Config = *ConfigT

// UnmarshalYAML ...
func (i Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return util.UnmarshalYAML(i, unmarshal)
}

// MarshalYAML ...
func (i Config) MarshalYAML() (interface{}, error) {
	return util.MarshalYAML(i)
}

// Reset ...
func (i Config) Reset() {
	i.Colorization = true
	i.Replace = make(map[string]string)
	i.Pattern = ""
	i.StartupLine.Reset()
	i.LineNo.Reset()
	i.UnknownLine.Reset()
	i.Prefix.Reset()
	i.Fields.Reset()
}

// FromMap ...
func (i Config) FromMap(m map[string]interface{}) error {
	var v interface{}

	v = util.ExtractFromMap(m, "colorization")
	if v != nil {
		i.Colorization = util.ToBool(v)
	}

	v = util.ExtractFromMap(m, "replace")
	if v != nil {
		i.Replace = v.(map[string]string)
	}

	v = util.ExtractFromMap(m, "pattern")
	if v != nil {
		i.Pattern = v.(string)
	}

	v = util.ExtractFromMap(m, "startup-line")
	if v != nil {
		if err := util.UnmashalYAMLAgain(v, &i.StartupLine); err != nil {
			return err
		}
	}

	v = util.ExtractFromMap(m, "line-no")
	if v != nil {
		if err := util.UnmashalYAMLAgain(v, &i.LineNo); err != nil {
			return err
		}
	}

	v = util.ExtractFromMap(m, "unknown-line")
	if v != nil {
		if err := util.UnmashalYAMLAgain(v, &i.UnknownLine); err != nil {
			return err
		}
	}

	v = util.ExtractFromMap(m, "prefix")
	if v != nil {
		if err := util.UnmashalYAMLAgain(v, &i.Prefix); err != nil {
			return err
		}
	}

	v = util.ExtractFromMap(m, "fields")
	if v != nil {
		if err := util.UnmashalYAMLAgain(v, &i.Fields); err != nil {
			return err
		}
	}

	return nil
}

// ToMap ...
func (i Config) ToMap() map[string]interface{} {
	r := make(map[string]interface{})
	r["replace"] = i.Replace
	r["pattern"] = i.Pattern
	r["startup-line"] = i.StartupLine.ToMap()
	r["line-no"] = i.LineNo.ToMap()
	r["unknown-line"] = i.UnknownLine.ToMap()
	r["prefix"] = i.Prefix.ToMap()
	r["fields"] = i.Fields.ToMap()
	return r
}

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

	dir, err := homedir.Dir()
	if err != nil {
		log.Printf("failed to get home dir: %v\n", err)
		return ""
	}
	return lookForConfigFile(dir)
}

// ConfigWithDefaultYamlFile ...
func ConfigWithDefaultYamlFile() Config {
	path := DetermineConfigFilePath()

	if len(path) == 0 {
		log.Println("config file not found, take default config")
		return ConfigWithYamlFile(config.DefaultYAML)
	}

	log.Printf("config file: %s\n", path)
	return ConfigWithYamlFile(path)
}

// ConfigWithYamlFile ...
func ConfigWithYamlFile(path string) Config {
	log.Printf("config file: %s\n", path)

	yamlText := string(util.ReadFile(path))
	return ConfigWithYaml(yamlText)
}

// ConfigWithYaml ...
func ConfigWithYaml(yamlText string) Config {
	r := &ConfigT{
		Replace: map[string]string{
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
