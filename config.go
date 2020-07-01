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

	Label  string
	Colors []color.Color
}

// NewColorConfig ...
func NewColorConfig(label string) ColorConfigT {
	var r ColorConfigT

	if err := yaml.UnmarshalStrict([]byte(label), &r); err != nil {
		panic(errors.Wrap(err, "failed to unmarshal color: "+label))
	}

	return r
}

// UnmarshalYAML ...
func (me ColorConfigT) UnmarshalYAML(unmarshal func(interface{}) error) error {
	if err := unmarshal(&me.Label); err != nil {
		return err
	}

	c, err := ColorsFromLabel(me.Label)
	if err != nil {
		return nil
	}
	me.Colors = c

	return nil
}

// MarshalYAML ...
func (me ColorConfigT) MarshalYAML() (interface{}, error) {
	return me.Label, nil
}

// OutputColorsConfigT ...
type OutputColorsConfigT struct {
	LineNo     ColorConfigT `yaml:"lineNo"`
	Timestamp  ColorConfigT `yaml:"timestamp"`
	Version    ColorConfigT `yaml:"version"`
	Message    ColorConfigT `yaml:"message"`
	Logger     ColorConfigT `yaml:"logger"`
	Thread     ColorConfigT `yaml:"thread"`
	StackTrace ColorConfigT `yaml:"stackTrace"`

	Debug ColorConfigT `yaml:"debug"`
	Info  ColorConfigT `yaml:"info"`
	Error ColorConfigT `yaml:"error"`
	Warn  ColorConfigT `yaml:"warn"`
	Trace ColorConfigT `yaml:"trace"`
	Fine  ColorConfigT `yaml:"fine"`
	Fatal ColorConfigT `yaml:"fatal"`

	Raw    ColorConfigT `yaml:"raw"`
	Others ColorConfigT `yaml:"others"`
}

// OutputColorsConfig ...
type OutputColorsConfig = *OutputColorsConfigT

// OutputConfigT ...
type OutputConfigT struct {
	Pattern string              `yaml:"pattern"`
	Colors  OutputColorsConfigT `yaml:"colors"`
}

// ConfigT ...
type ConfigT struct {
	Output OutputConfigT `yaml:"output"`
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

// LoadConfig ...
func LoadConfig() Config {
	color.Red.Println("Simple to use color")

	path := DetermineConfigFilePath()
	log.Printf("Config file: %s\n", path)

	r := &ConfigT{
		Output: OutputConfigT{
			Pattern: "",
			Colors: OutputColorsConfigT{
				LineNo:     NewColorConfig("FgWhite"),
				Timestamp:  NewColorConfig("FgWhite"),
				Version:    NewColorConfig("FgWhite"),
				Message:    NewColorConfig("FgWhite"),
				Logger:     NewColorConfig("FgWhite"),
				Thread:     NewColorConfig("FgWhite"),
				StackTrace: NewColorConfig("FgWhite"),

				Debug: NewColorConfig("FgBlue,OpBold"),
				Info:  NewColorConfig("FgWhite,OpBold"),
				Error: NewColorConfig("FgRed,OpBold"),
				Warn:  NewColorConfig("FgYellow,OpBold"),
				Trace: NewColorConfig("FgGreen,OpBold"),
				Fine:  NewColorConfig("FgCyan,OpBold"),
				Fatal: NewColorConfig("FgRed,OpBold"),

				Raw:    NewColorConfig("FgWhite"),
				Others: NewColorConfig("FgWhite"),
			},
		},
	}

	if len(path) != 0 {
		raw := ReadFile(path)

		if err := yaml.UnmarshalStrict(raw, &r); err != nil {
			panic(errors.Wrap(err, "failed to unmarshal config file: "+path))
		}
	}

	return r
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
