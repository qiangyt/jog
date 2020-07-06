package config

import (
	"strings"

	"gopkg.in/yaml.v2"
)

// LoggerNameCompressT ...
type LoggerNameCompressT struct {
	yaml.Unmarshaler
	yaml.Marshaler

	Enabled   bool
	Separator string
}

// LoggerNameCompress ..
type LoggerNameCompress = *LoggerNameCompressT

// Set ...
func (i LoggerNameCompress) Set(valuesText string) {
	i.Separator = "."

	t := strings.ToLower(valuesText)
	if "true" == t || "t" == t || "yes" == t || "y" == t {
		i.Enabled = true
	} else if "false" == t || "f" == t || "no" == t || "n" == t {
		i.Enabled = false
	} else {
		i.Enabled = true
		i.Separator = t
	}
}

// Reset ...
func (i LoggerNameCompress) Reset() {
	i.Separator = ""
	i.Enabled = false
}

// Value ...
func (i LoggerNameCompress) Value() interface{} {
	if i.Enabled {
		return i.Separator
	}
	return false
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
	return i.Value(), nil
}

var _loggerNameCache = make(map[string]string)

func (i LoggerNameCompress) compressIfEnabled(loggerName string) string {
	if !i.Enabled {
		return loggerName
	}

	if existingOne, ok := _loggerNameCache[loggerName]; ok {
		return existingOne
	}

	pkgList := strings.Split(loggerName, i.Separator)
	indexLast := len(pkgList) - 1
	for index, pkg := range pkgList[:indexLast] {
		pkgList[index] = string([]byte(pkg)[0])
	}

	r := strings.Join(pkgList, i.Separator)
	_loggerNameCache[loggerName] = r
	return r
}
