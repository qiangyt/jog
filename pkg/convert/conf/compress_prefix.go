package conf

import (
	"fmt"
	"strings"

	"github.com/gookit/goutil/strutil"
	_util "github.com/qiangyt/jog/pkg/util"
	"gopkg.in/yaml.v2"
)

// CompressPrefixAction_ ...
type CompressPrefixAction_ int

const (
	// CompressPrefixAction_RemoveNonFirstLetter ...
	CompressPrefixAction_RemoveNonFirstLetter CompressPrefixAction_ = iota

	// CompressPrefixAction_Remove ...
	CompressPrefixAction_Remove

	// CompressPrefixAction_Default ...
	CompressPrefixAction_Default = CompressPrefixAction_RemoveNonFirstLetter
)

// Format ...
func (i CompressPrefixAction_) String() string {
	if i == CompressPrefixAction_RemoveNonFirstLetter {
		return "remove-non-first-letter"
	}
	if i == CompressPrefixAction_Remove {
		return "remove"
	}

	return ""
}

// ParseCompressPrefixAction ...
func ParseCompressPrefixAction(text string) CompressPrefixAction_ {
	if "remove-non-first-letter" == text {
		return CompressPrefixAction_RemoveNonFirstLetter
	}
	if "remove" == text {
		return CompressPrefixAction_Remove
	}

	panic(fmt.Errorf("unknown CompressPrefixAction_ text: %v", text))
}

// CompressPrefixT ...
type CompressPrefixT struct {
	yaml.Unmarshaler
	yaml.Marshaler

	Enabled    bool
	Separators StringSet
	WhiteList  StringSet
	Action     CompressPrefixAction_
}

// CompressPrefix ..
type CompressPrefix = *CompressPrefixT

// UnmarshalYAML ...
func (i CompressPrefix) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return _util.DynObject4YAML(i, unmarshal)
}

// MarshalYAML ...
func (i CompressPrefix) MarshalYAML() (interface{}, error) {
	return _util.DynObject2YAML(i)
}

// FromMap ...
func (i CompressPrefix) FromMap(m map[string]interface{}) error {

	enabledV := _util.ExtractFromMap(m, "enabled")
	if enabledV != nil {
		i.Enabled = _util.ToBool(enabledV)
	}

	separatorsV := _util.ExtractFromMap(m, "separators")
	if separatorsV != nil {
		i.Separators.Parse(separatorsV)
	}

	whiteListV := _util.ExtractFromMap(m, "white-list")
	if whiteListV != nil {
		i.WhiteList.Parse(whiteListV)
	}

	actionV := _util.ExtractFromMap(m, "action")
	if actionV != nil {
		i.Action = ParseCompressPrefixAction(strutil.MustString(actionV))
	}

	return nil
}

// ToMap ...
func (i CompressPrefix) ToMap() map[string]interface{} {
	r := make(map[string]interface{})
	r["enabled"] = i.Enabled
	r["separators"] = i.Separators.String()
	r["white-list"] = i.WhiteList.String()
	r["action"] = i.Action.String()
	return r
}

// Reset ...
func (i CompressPrefix) Reset() {
	i.Enabled = false
	i.Separators = &StringSetT{}
	i.WhiteList = &StringSetT{}
	i.Action = CompressPrefixAction_Default
}

// TODO: this 2 caches are not thread-safe, should be moved to a context
var _compressCache4RemoveNonFirstLetter = make(map[string]string)
var _compressCache4Remove = make(map[string]string)

func (i CompressPrefix) detectSeparator(text string) (string, []string) {
	for separator := range i.Separators.ValueMap {
		separated := strings.Split(text, separator)
		if len(separated) > 1 {
			return separator, separated
		}
	}

	if i.Separators.CaseSensitive == false {
		for separator := range i.Separators.LowercasedValueMap {
			separated := strings.Split(text, separator)
			if len(separated) > 1 {
				return separator, separated
			}
		}
		for separator := range i.Separators.UppercasedValueMap {
			separated := strings.Split(text, separator)
			if len(separated) > 1 {
				return separator, separated
			}
		}
	}

	return "", []string{text}
}

// Compress ...
func (i CompressPrefix) Compress(text string) string {
	if i.Action == CompressPrefixAction_RemoveNonFirstLetter {
		if existingOne, ok := _compressCache4RemoveNonFirstLetter[text]; ok {
			return existingOne
		}

		var r string
		if i.WhiteList.ContainsPrefixOf(text) {
			r = text
		} else {
			separator, separated := i.detectSeparator(text)

			if len(separated) > 1 {
				indexOfLast := len(separated) - 1
				for index, item := range separated[:indexOfLast] {
					separated[index] = string([]byte(item)[0])
				}

				r = strings.Join(separated, separator)
			} else {
				r = text
			}
		}

		_compressCache4RemoveNonFirstLetter[text] = r
		return r
	}

	if i.Action == CompressPrefixAction_Remove {
		if existingOne, ok := _compressCache4Remove[text]; ok {
			return existingOne
		}

		var r string
		if i.WhiteList.ContainsPrefixOf(text) {
			r = text
		} else {
			_, separated := i.detectSeparator(text)

			if len(separated) > 1 {
				r = separated[len(separated)-1]
			} else {
				r = text
			}
		}
		_compressCache4Remove[text] = r
		return r
	}

	return text
}
