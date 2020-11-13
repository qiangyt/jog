package config

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/qiangyt/jog/static"
	"github.com/qiangyt/jog/static/grok_vjeantet"
	"github.com/qiangyt/jog/util"
	"gopkg.in/yaml.v2"
)

const (
	jogHomeDir = "~/.jog"
)

// JogHomeDir ...
func JogHomeDir(expand bool, children ...string) string {
	var r string

	if !expand {
		r = jogHomeDir
	} else {
		r = util.ExpandPath(jogHomeDir)
		util.MkdirAll(r)
	}

	return filepath.Join(r, filepath.Join(children...))
}

// ConfigurationT ...
type ConfigurationT struct {
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

// Configuration ...
type Configuration = *ConfigurationT

// UnmarshalYAML ...
func (i Configuration) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return UnmarshalYAML(i, unmarshal)
}

// MarshalYAML ...
func (i Configuration) MarshalYAML() (interface{}, error) {
	return MarshalYAML(i)
}

// Init ...
func (i Configuration) Init(cfg Configuration) {
	if cfg != nil {
		panic(fmt.Errorf("root configure initialization"))
	}

	i.StartupLine.Init(i)
	i.LineNo.Init(i)
	i.UnknownLine.Init(i)
	i.Prefix.Init(i)
	i.Fields.Init(i)

	levelField := i.Fields.Standards["level"]
	if levelField != nil {
		if !levelField.IsEnum() {
			panic(fmt.Errorf("invalid configuration: field 'level' must be enum"))
		}
	}
	i.LevelField = levelField

	timestampField := i.Fields.Standards["timestamp"]
	if timestampField != nil {
		if timestampField.Type == FieldTypeAuto {
			timestampField.Type = FieldTypeTime
		} else if timestampField.Type != FieldTypeTime {
			panic(fmt.Errorf("invalid configuration: type of field 'timestamp' must be 'time' or 'auto'"))
		}
	}
	i.TimestampField = timestampField

	i.Grok.Init(i)
}

// Reset ...
func (i Configuration) Reset() {
	i.Colorization = true
	i.Replace = make(map[string]string)
	i.Pattern = ""
	i.HasOthersFieldInPattern = false
	i.fieldsInPattern = make(map[string]bool)
	i.StartupLine.Reset()
	i.LineNo.Reset()
	i.UnknownLine.Reset()
	i.Prefix.Reset()
	i.Fields.Reset()
	i.LevelField = nil
	i.Grok.Reset()
}

// HasFieldInPattern ...
func (i Configuration) HasFieldInPattern(fieldName string) bool {
	r, contains := i.fieldsInPattern[fieldName]
	if contains {
		return r
	}

	r = strings.Contains(i.Pattern, "${"+fieldName+"}")
	i.fieldsInPattern[fieldName] = r
	return r
}

// FromMap ...
func (i Configuration) FromMap(m map[string]interface{}) error {
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
func (i Configuration) ToMap() map[string]interface{} {
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

type grokPatternItem struct {
	Name string
	Expr string
}

func parseGrokPatterns(patternsText string) []grokPatternItem {
	r := make([]grokPatternItem, 0)

	buf := bytes.NewBufferString(patternsText)

	scanner := bufio.NewScanner(buf)
	for scanner.Scan() {
		l := scanner.Text()
		l = strings.TrimSpace(l)
		if len(l) > 0 && l[0] != '#' {
			nameAndExpr := strings.SplitN(l, " ", 2)

			p := grokPatternItem{}
			p.Name = nameAndExpr[0]
			p.Expr = nameAndExpr[1]

			r = append(r, p)
		}
	}

	return r
}

func loadAllGrokPatterns() []grokPatternItem {
	r := []grokPatternItem{}

	r = append(r, parseGrokPatterns(grok_vjeantet.Aws)...)
	r = append(r, parseGrokPatterns(grok_vjeantet.Bacula)...)
	r = append(r, parseGrokPatterns(grok_vjeantet.Bro)...)
	r = append(r, parseGrokPatterns(grok_vjeantet.Exim)...)
	r = append(r, parseGrokPatterns(grok_vjeantet.Firewalls)...)
	r = append(r, parseGrokPatterns(grok_vjeantet.Grok_patterns)...)
	r = append(r, parseGrokPatterns(grok_vjeantet.Haproxy)...)
	r = append(r, parseGrokPatterns(grok_vjeantet.Java)...)
	r = append(r, parseGrokPatterns(grok_vjeantet.Junos)...)
	r = append(r, parseGrokPatterns(grok_vjeantet.Linux_syslog)...)
	r = append(r, parseGrokPatterns(grok_vjeantet.Mcollective_patterns)...)
	r = append(r, parseGrokPatterns(grok_vjeantet.Mcollective)...)
	r = append(r, parseGrokPatterns(grok_vjeantet.Mongodb)...)
	r = append(r, parseGrokPatterns(grok_vjeantet.Nagios)...)
	r = append(r, parseGrokPatterns(grok_vjeantet.Postgresql)...)
	r = append(r, parseGrokPatterns(grok_vjeantet.Rails)...)
	r = append(r, parseGrokPatterns(grok_vjeantet.Redis)...)
	r = append(r, parseGrokPatterns(grok_vjeantet.Ruby)...)

	return r
}

func mergeGrokPatterns(allPatterns map[string]grokPatternItem, patternsText string) {
	newPatterns := parseGrokPatterns(patternsText)
	for _, pattern := range newPatterns {
		name := pattern.Name
		if existingOne, alreadyExists := allPatterns[name]; alreadyExists == true {
			panic(fmt.Errorf("duplicated grok pattern. name: %s. existing: %s. duplicated: %s", name, existingOne.Expr, pattern.Expr))
		}
		allPatterns[name] = pattern
	}
}

// BuildDefaultConfigurationYAML ...
func BuildDefaultConfigurationYAML() string {
	tmpl, err := template.New("default configuration YAML").Parse(static.DefaultConfiguration_yml)
	if err != nil {
		panic(errors.Wrap(err, "failed to parse default configuration YAML as template"))
	}

	grokPatterns := loadAllGrokPatterns()

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, map[string]interface{}{"grokPatterns": grokPatterns})
	if err != nil {
		panic(errors.Wrap(err, "failed to execute default configuration YAML as template"))
	}
	return buf.String()
}

// WithDefaultYamlFile ...
func WithDefaultYamlFile() Configuration {
	path := DetermineConfigFilePath()

	if len(path) == 0 {
		log.Println("config file not found, take default one")
		return WithYaml(BuildDefaultConfigurationYAML())
	}

	log.Printf("config file: %s\n", path)
	return WithYamlFile(path)
}

// WithYamlFile ...
func WithYamlFile(path string) Configuration {
	log.Printf("config file: %s\n", path)

	yamlText := string(util.ReadFile(path))
	return WithYaml(yamlText)
}

// WithYaml ...
func WithYaml(yamlText string) Configuration {
	r := &ConfigurationT{
		Replace: map[string]string{
			"\\\"": "\"",
			"\\'":  "'",
			"\\\n": "\n",
			"\\\r": "",
			"\\\t": "\t",
		},
		Pattern:     "",
		StartupLine: &StartupLineT{},
		LineNo:      &ElementT{},
		UnknownLine: &ElementT{},
		Prefix:      &PrefixT{},
		Fields:      &FieldMapT{},
		Grok:        &GrokT{},
	}

	if err := yaml.Unmarshal([]byte(yamlText), &r); err != nil {
		panic(errors.Wrap(err, "failed to unmarshal yaml: \n"+yamlText))
	}

	r.Init(nil)

	return r
}
