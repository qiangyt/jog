package config

import (
	"strings"

	"gopkg.in/yaml.v2"
)

// LoggerNameCompressT ...
type LoggerNameCompressT struct {
	yaml.Unmarshaler
	yaml.Marshaler

	enabled   bool
	separator string
}

// LoggerNameCompress ..
type LoggerNameCompress = *LoggerNameCompressT

// Set ...
func (i LoggerNameCompress) Set(valuesText string) {
	i.separator = "."

	t := strings.ToLower(valuesText)
	if "true" == t || "t" == t || "yes" == t || "y" == t {
		i.enabled = true
	} else if "false" == t || "f" == t || "no" == t || "n" == t {
		i.enabled = false
	} else {
		i.enabled = true
		i.separator = t
	}
}

// Reset ...
func (i LoggerNameCompress) Reset() {
	i.separator = ""
	i.enabled = false
}

func (i LoggerNameCompress) String() string {
	if i.enabled {
		return i.separator
	}
	return "false"
}

// UnmarshalYAML ...
func (i LoggerNameCompress) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var valuesText string
	err := unmarshal(&valuesText)
	if err != nil {
		return err
	}

	i.Set(valuesText)

	return nil
}

// MarshalYAML ...
func (i LoggerNameCompress) MarshalYAML() (interface{}, error) {
	return i.String(), nil
}

var _loggerNameCache = make(map[string]string)

func (i LoggerNameCompress) compressIfEnabled(loggerName string) string {
	if !i.enabled {
		return loggerName
	}

	if existingOne, ok := _loggerNameCache[loggerName]; ok {
		return existingOne
	}

	pkgList := strings.Split(loggerName, i.separator)
	indexLast := len(pkgList) - 1
	for index, pkg := range pkgList[:indexLast] {
		pkgList[index] = string([]byte(pkg)[0])
	}

	r := strings.Join(pkgList, i.separator)
	_loggerNameCache[loggerName] = r
	return r
}
