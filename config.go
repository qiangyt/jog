package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"runtime"

	"github.com/gookit/color"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// OSType ...
type OSType int

const (
	// AllOSType ...
	AllOSType OSType = 0

	// Windows ...
	Windows OSType = 1

	// Darwin ...
	Darwin OSType = 2

	// Linux ...
	Linux OSType = 3
)

// ParseOSType ...
func ParseOSType(s string) OSType {
	if "all" == s {
		return AllOSType
	}
	if "windows" == s {
		return Windows
	}
	if "darwin" == s {
		return Darwin
	}
	if "linux" == s {
		return Linux
	}
	panic(fmt.Errorf("unknown OS type: '%s'", s))
}

// BuildOSType ...
func BuildOSType(i int) OSType {
	if int(AllOSType) == i {
		return AllOSType
	}
	if int(Windows) == i {
		return Windows
	}
	if int(Darwin) == i {
		return Darwin
	}
	if int(Linux) == i {
		return Linux
	}
	panic(fmt.Errorf("unknown OS type: '%v'", i))
}

// DefaultOSType ...
func DefaultOSType() OSType {
	return ParseOSType(runtime.GOOS)
}

// ColorConfigT ...
type ColorConfigT struct {
	yaml.Unmarshaler
	yaml.Marshaler

	Label string
	Style color.Style
}

// ColorConfig ...
type ColorConfig = *ColorConfigT

// NewColorConfig ...
func NewColorConfig(label string) ColorConfig {
	var r ColorConfigT

	if err := yaml.UnmarshalStrict([]byte(label), &r); err != nil {
		panic(errors.Wrap(err, "failed to unmarshal color: "+label))
	}

	return &r
}

// UnmarshalYAML ...
func (me ColorConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	err := unmarshal(&me.Label)
	if err != nil {
		return err
	}

	me.Style, err = ColorsFromLabel(me.Label)
	return err
}

// MarshalYAML ...
func (me ColorConfig) MarshalYAML() (interface{}, error) {
	return me.Label, nil
}

// Render ...
func (me ColorConfig) Render(a ...interface{}) string {
	return me.Style.Render(a...)
}

// Sprint is alias of the 'Render'
func (me ColorConfig) Sprint(a ...interface{}) string {
	return me.Style.Sprint(a...)
}

// Sprintf format and render message.
func (me ColorConfig) Sprintf(format string, a ...interface{}) string {
	return me.Style.Sprintf(format, a...)
}

// Print render and Print text
func (me ColorConfig) Print(a ...interface{}) {
	me.Style.Print(a...)
}

// Printf render and print text
func (me ColorConfig) Printf(format string, a ...interface{}) {
	me.Style.Printf(format, a...)
}

// Println render and print text line
func (me ColorConfig) Println(a ...interface{}) {
	me.Style.Println(a...)
}

// OutputColorsConfigT ...
type OutputColorsConfigT struct {
	LineNo      ColorConfig `yaml:"line-no"`
	Timestamp   ColorConfig
	Version     ColorConfig
	Message     ColorConfig
	Logger      ColorConfig
	Thread      ColorConfig
	StackTrace  ColorConfig `yaml:"stack-trace"`
	StartedLine ColorConfig `yaml:"started-line"`

	Debug ColorConfig
	Info  ColorConfig
	Error ColorConfig
	Warn  ColorConfig
	Trace ColorConfig
	Fine  ColorConfig
	Fatal ColorConfig

	Raw             ColorConfig
	OthersName      ColorConfig `yaml:"others-name"`
	OthersSeparator ColorConfig `yaml:"others-separator"`
	OthersValue     ColorConfig `yaml:"others-value"`
}

// OutputColorsConfig ...
type OutputColorsConfig = *OutputColorsConfigT

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

// DefaultConfigYAML ...
const DefaultConfigYAML = `
output:
  pattern: "${timestamp} ${level} <${thread}> ${logger}: ${message} ${others} ${stacktrace}"
  compress-logger-name: true
  colors:
    line-no: FgDefault
    timestamp: FgDefault
    version: FgDefault
    message: FgDefault
    logger: FgDefault
    thread: FgDefault
    stack-trace: FgDefault
    started-line: FgGreen, OpBold

    debug: FgBlue,OpBold
    info: FgBlue,OpBold
    error: FgRed,OpBold
    warn: FgYellow,OpBold
    trace: FgBlue,OpBold
    fine: FgCyan,OpBold
    fatal: FgRed,OpBold

    raw: FgDefault
    others-name: FgDefault,OpBold
    others-separator: FgDefault
    others-value: FgDefault

  started-line: Started Application in
`

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
		return ConfigWithYaml(DefaultConfigYAML)
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

func getConfigString(prefix string, m map[string]interface{}, key string) (bool, string) {
	v, has := m[key]
	if !has {
		return false, ""
	}

	r, is := v.(string)
	if !is {
		panic(errors.Errorf("%s.%s is expected to be a string, but this is a %v", prefix, key, reflect.TypeOf(v)))
	}
	return true, r
}

// GetConfigStringD ...
func GetConfigStringD(prefix string, m map[string]interface{}, key string, def string) string {
	has, r := getConfigString(prefix, m, key)
	if !has {
		return def
	}
	return r
}

// GetConfigStringP ...
func GetConfigStringP(prefix string, m map[string]interface{}, key string) string {
	has, r := getConfigString(prefix, m, key)
	if !has {
		panic(fmt.Errorf("%s missing %s", prefix, key))
	}
	return r
}

func getConfigBool(prefix string, m map[string]interface{}, key string) (bool, bool) {
	v, has := m[key]
	if !has {
		return false, false
	}

	r, is := v.(bool)
	if !is {
		panic(errors.Errorf("%s.%s is expected to be a bool, but this is a %v", prefix, key, reflect.TypeOf(v)))
	}
	return true, r
}

// GetConfigBoolD ...
func GetConfigBoolD(prefix string, m map[string]interface{}, key string, def bool) bool {
	has, r := getConfigBool(prefix, m, key)
	if !has {
		return def
	}
	return r
}

// GetConfigBoolP ...
func GetConfigBoolP(prefix string, m map[string]interface{}, key string) bool {
	has, r := getConfigBool(prefix, m, key)
	if !has {
		panic(fmt.Errorf("%s missing %s", prefix, key))
	}
	return r
}

func getConfigObject(prefix string, m map[string]interface{}, key string) (bool, map[string]interface{}) {
	v, has := m[key]
	if !has {
		return false, nil
	}

	r, is := v.(map[interface{}]interface{})
	if !is {
		panic(errors.Errorf("%s.%s is expected to be an object, but this is a %v", prefix, key, reflect.TypeOf(v)))
	}

	return true, MapToObject(r)
}

// GetConfigObjectD ...
func GetConfigObjectD(prefix string, m map[string]interface{}, key string, def map[string]interface{}) map[string]interface{} {
	has, r := getConfigObject(prefix, m, key)
	if !has {
		return def
	}
	return r
}

// GetConfigObjectP ...
func GetConfigObjectP(prefix string, m map[string]interface{}, key string) map[string]interface{} {
	has, r := getConfigObject(prefix, m, key)
	if !has {
		panic(fmt.Errorf("%s missing %s", prefix, key))
	}
	return r
}

func getConfigObjectArray(prefix string, m map[string]interface{}, key string) (bool, []interface{}) {
	v, has := m[key]
	if !has {
		return false, nil
	}

	r, is := v.([]interface{})
	if !is {
		panic(errors.Errorf("%s.%s is expected to be an array of object, but this is a %v", prefix, key, reflect.TypeOf(v)))
	}
	return true, r
}

// GetConfigObjectArrayD ...
func GetConfigObjectArrayD(prefix string, m map[string]interface{}, key string, def []interface{}) []interface{} {
	has, r := getConfigObjectArray(prefix, m, key)
	if !has {
		return def
	}
	return r
}

// GetConfigObjectArrayP ...
func GetConfigObjectArrayP(prefix string, m map[string]interface{}, key string) []interface{} {
	has, r := getConfigObjectArray(prefix, m, key)
	if !has {
		panic(fmt.Errorf("%s missing %s", prefix, key))
	}
	return r
}

func getConfigStringArray(prefix string, m map[string]interface{}, key string) (bool, []string) {
	v, has := m[key]
	if !has {
		return false, nil
	}

	r, is := v.([]interface{})
	if !is {
		panic(errors.Errorf("%s.%s is expected to be an array of string, but this is a %v", prefix, key, reflect.TypeOf(v)))
	}
	return true, NormalizeStringArray(r)
}

// GetConfigStringArrayD ...
func GetConfigStringArrayD(prefix string, m map[string]interface{}, key string, def []string) []string {
	has, r := getConfigStringArray(prefix, m, key)
	if !has {
		return def
	}
	return r
}

// GetConfigStringArrayP ...
func GetConfigStringArrayP(prefix string, m map[string]interface{}, key string) []string {
	has, r := getConfigStringArray(prefix, m, key)
	if !has {
		panic(fmt.Errorf("%s missing %s", prefix, key))
	}
	return r
}

func getConfigInt(prefix string, m map[string]interface{}, key string) (bool, int) {
	v, has := m[key]
	if !has {
		return false, 0
	}

	r, is := v.(int)
	if !is {
		panic(errors.Errorf("%s.%s is expected to be a int, but this is a %v", prefix, key, reflect.TypeOf(v)))
	}
	return true, r
}

// GetConfigIntD ...
func GetConfigIntD(prefix string, m map[string]interface{}, key string, def int) int {
	has, r := getConfigInt(prefix, m, key)
	if !has {
		return def
	}
	return r
}

// GetConfigIntP ...
func GetConfigIntP(prefix string, m map[string]interface{}, key string) int {
	has, r := getConfigInt(prefix, m, key)
	if !has {
		panic(fmt.Errorf("%s missing %s", prefix, key))
	}
	return r
}

func getConfigOSType(prefix string, m map[string]interface{}, key string) (bool, OSType) {
	v, has := m[key]
	if !has {
		return false, Windows
	}

	r, is := v.(string)
	if !is {
		i, is := v.(int)
		if !is {
			panic(errors.Errorf("%s.%s is expected to be a string, but this is a %v", prefix, key, reflect.TypeOf(v)))
		}
		return true, BuildOSType(i)
	}
	return true, ParseOSType(r)
}

// GetConfigOSTypeD ...
func GetConfigOSTypeD(prefix string, m map[string]interface{}, key string, def OSType) OSType {
	has, r := getConfigOSType(prefix, m, key)
	if !has {
		return def
	}
	return r
}

// GetConfigOSTypeP ...
func GetConfigOSTypeP(prefix string, m map[string]interface{}, key string) OSType {
	has, r := getConfigOSType(prefix, m, key)
	if !has {
		panic(fmt.Errorf("%s missing %s", prefix, key))
	}
	return r
}

func getConfigMap(prefix string, m map[string]interface{}, key string) (bool, map[string]interface{}) {
	v, has := m[key]
	if !has {
		return false, nil
	}

	r, is := v.(map[string]interface{})
	if !is {
		panic(errors.Errorf("%s.%s is expected to be a map/object, but this is a %v", prefix, key, reflect.TypeOf(v)))
	}
	return true, r
}

// GetConfigMapD ...
func GetConfigMapD(prefix string, m map[string]interface{}, key string, def map[string]interface{}) map[string]interface{} {
	has, r := getConfigMap(prefix, m, key)
	if !has {
		return def
	}
	return r
}

// GetConfigMapP ...
func GetConfigMapP(prefix string, m map[string]interface{}, key string) map[string]interface{} {
	has, r := getConfigMap(prefix, m, key)
	if !has {
		panic(fmt.Errorf("%s missing %s", prefix, key))
	}
	return r
}
