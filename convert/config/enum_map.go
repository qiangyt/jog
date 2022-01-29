package config

import (
	"fmt"
	"strings"

	"github.com/gookit/goutil/strutil"
	"github.com/pkg/errors"
	"github.com/qiangyt/jog/util"
)

// EnumMapT ...
type EnumMapT struct {
	CaseSensitive bool `yaml:"case-sensitive"`
	Default       string
	values        map[string]Enum // normalized name -> enum
	allMap        map[string]Enum // normalized name + aliases -> enum
}

// EnumMap ...
type EnumMap = *EnumMapT

// Reset ...
func (i EnumMap) Reset() {
	i.CaseSensitive = false
	i.Default = ""
	i.values = make(map[string]Enum)
	i.allMap = make(map[string]Enum)
}

// UnmarshalYAML ...
func (i EnumMap) UnmarshalYAML(unmarshal func(interface{}) error) error {
	e := UnmarshalYAML(i, unmarshal)
	return e
}

// MarshalYAML ...
func (i EnumMap) MarshalYAML() (interface{}, error) {
	return MarshalYAML(i)
}

// Init ...
func (i EnumMap) Init(cfg Configuration) {

}

// IsEmpty ...
func (i EnumMap) IsEmpty() bool {
	return len(i.values) == 0
}

// GetEnum ...
func (i EnumMap) GetEnum(value string) Enum {
	if !i.CaseSensitive {
		value = strings.ToLower(value)
	}

	r, has := i.allMap[value]
	if has {
		return r
	}
	return i.allMap[i.Default]
}

// FromMap ...
func (i EnumMap) FromMap(m map[string]interface{}) error {
	caseSensitiveV := util.ExtractFromMap(m, "case-sensitive")
	if caseSensitiveV != nil {
		i.CaseSensitive = util.ToBool(caseSensitiveV)
	}

	defaultV := util.ExtractFromMap(m, "default")
	if defaultV != nil {
		i.Default = strutil.MustString(defaultV)
	}

	for k, v := range m {
		var ev Enum
		if err := util.UnmashalYAMLAgain(v, &ev); err != nil {
			return err
		}

		ev.Name = k
		i.values[k] = ev

		i.allMap[k] = ev
		if !i.CaseSensitive {
			i.allMap[strings.ToLower(k)] = ev
		}

		for aliasName := range ev.Alias.Values {
			i.allMap[aliasName] = ev
			if !i.CaseSensitive {
				i.allMap[strings.ToLower(aliasName)] = ev
			}
		}

		delete(m, k)

		if len(i.Default) == 0 {
			i.Default = k
		}
	}

	if len(i.Default) == 0 {
		return errors.New("default enum not specified")
	}
	if _, defaultMatches := i.allMap[i.Default]; !defaultMatches {
		return fmt.Errorf("invalid default enum name: %s", i.Default)
	}

	return nil
}

// ToMap ...
func (i EnumMap) ToMap() map[string]interface{} {
	r := make(map[string]interface{})
	r["case-sensitive"] = i.CaseSensitive
	r["default"] = i.Default

	for k, v := range i.values {
		r[k] = v.ToMap()
	}

	return r
}
