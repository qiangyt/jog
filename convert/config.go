package convert

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/qiangyt/jog/util"
	"gopkg.in/yaml.v2"
)

const (
	DefaultConfigFile = "jog.convert.yaml"
)

var _jogConvertYamlResource util.Resource

func init() {
	_jogConvertYamlResource = util.NewResource("/" + DefaultConfigFile)
}

// ConfigT ...
type ConfigT struct {
	// TODO: configurable
	Colorization            bool
	Replace                 map[string]string
	Pattern                 string
	fieldsInPattern         map[string]bool
	HasOthersFieldInPattern bool
	StartupLine             StartupLine `yaml:"startup-line"`
	LineNo                  Element     `yaml:"line-no"`
	UnknownLine             Element     `yaml:"unknown-line"`
	Prefix                  Prefix
	Fields                  FieldMap
	LevelField              Field
	TimestampField          Field
	Grok                    Grok
}

// Config ...
type Config = *ConfigT

// UnmarshalYAML ...
func (i Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return util.DynObject4YAML(i, unmarshal)
}

// MarshalYAML ...
func (i Config) MarshalYAML() (interface{}, error) {
	return util.DynObject2YAML(i)
}

// Init ...
func (i Config) Init(cfg Config) {
	/*i.StartupLine.Reset()
	i.LineNo.Reset()
	i.UnknownLine.Reset()
	i.Prefix.Reset()
	i.Fields.Reset()
	*/
	levelField := i.Fields.Standards["level"]
	if levelField != nil {
		if !levelField.IsEnum() {
			panic(fmt.Errorf("invalid configuration: field 'level' must be enum"))
		}
	}
	i.LevelField = levelField

	timestampField := i.Fields.Standards["timestamp"]
	if timestampField != nil {
		if timestampField.Type == FieldType_Auto {
			timestampField.Type = FieldType_Time
		} else if timestampField.Type != FieldType_Time {
			panic(fmt.Errorf("invalid configuration: type of field 'timestamp' must be 'time' or 'auto'"))
		}
	}
	i.TimestampField = timestampField

	util.InitDefaultGrokLibraryDir()
	i.Grok.Init(i)
}

// Reset ...
func (i Config) Reset() {
	i.Colorization = true

	i.Replace = map[string]string{
		"\\\"": "\"",
		"\\'":  "'",
		"\\\n": "\n",
		"\\\r": "",
		"\\\t": "\t",
	}

	i.Pattern = ""
	i.HasOthersFieldInPattern = false
	i.fieldsInPattern = make(map[string]bool)

	i.StartupLine = &StartupLineT{}
	i.StartupLine.Reset()

	i.LineNo = &ElementT{}
	i.LineNo.Reset()

	i.UnknownLine = &ElementT{}
	i.UnknownLine.Reset()

	i.Prefix = &PrefixT{}
	i.Prefix.Reset()

	i.Fields = &FieldMapT{}
	i.Fields.Reset()

	i.LevelField = nil

	i.Grok = &GrokT{}
	i.Grok.Reset()
}

// HasFieldInPattern ...
func (i Config) HasFieldInPattern(fieldName string) bool {
	r, contains := i.fieldsInPattern[fieldName]
	if contains {
		return r
	}

	r = strings.Contains(i.Pattern, "${"+fieldName+"}")
	i.fieldsInPattern[fieldName] = r
	return r
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
		i.HasOthersFieldInPattern = i.HasFieldInPattern("others")
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

	v = util.ExtractFromMap(m, "grok")
	if v != nil {
		if err := util.UnmashalYAMLAgain(v, &i.Grok); err != nil {
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
	r["grok"] = i.Grok.ToMap()
	return r
}

func LookForConfigFile(ctx ConvertContext, dir string) string {
	ctx.LogInfo("looking for config files", "dir", dir)
	r := filepath.Join(dir, DefaultConfigFile)
	if util.FileExists(r) {
		return r
	}
	r = filepath.Join(dir, DefaultConfigFile)
	if util.FileExists(r) {
		return r
	}
	return ""
}

// determineConfigFilePath return (file path)
func determineConfigFilePath(ctx ConvertContext) string {
	exeDir := util.ExeDirectory()
	r := LookForConfigFile(ctx, exeDir)
	if len(r) != 0 {
		return r
	}

	homeDir, err := homedir.Dir()
	if err != nil {
		ctx.LogInfo("failed to get home dir", "err", err)
	} else {
		r = LookForConfigFile(ctx, homeDir)
	}
	return r
}

// BuildDefaultConfigYAML ...
func BuildDefaultConfigYAML() string {
	yaml := util.NewResource(filepath.Join("/", DefaultConfigFile)).ReadString()

	tmpl, err := template.New("default configuration YAML").Parse(string(yaml))
	if err != nil {
		panic(errors.Wrap(err, "failed to parse default configuration YAML as template"))
	}

	grokPatterns := util.LoadAllGrokPatterns()

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, map[string]interface{}{"grokPatterns": grokPatterns})
	if err != nil {
		panic(errors.Wrap(err, "failed to execute default configuration YAML as template"))
	}
	return buf.String()
}

// NewConfigWithDefaultYamlFile ...
func NewConfigWithDefaultYamlFile(ctx ConvertContext) Config {
	configFilePath := determineConfigFilePath(ctx)

	if len(configFilePath) == 0 {
		ctx.LogInfo("config file not found, take default one")
		return NewConfigWithYaml(BuildDefaultConfigYAML())
	}

	ctx.LogInfo("config file", "path", configFilePath)
	return NewConfigWithYamlFile(ctx, configFilePath)
}

// NewConfigWithYamlFile ...
func NewConfigWithYamlFile(ctx ConvertContext, path string) Config {
	ctx.LogInfo("config file", "path", path)

	yamlText := string(util.ReadFile(path))
	return NewConfigWithYaml(yamlText)
}

// NewConfigWithYaml ...
func NewConfigWithYaml(yamlText string) Config {
	r := &ConfigT{}
	r.Reset()

	if err := yaml.Unmarshal([]byte(yamlText), &r); err != nil {
		panic(errors.Wrap(err, "failed to unmarshal yaml: \n"+yamlText))
	}

	r.Init(nil)

	return r
}
