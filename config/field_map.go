package config

import (
	"fmt"
	"strings"

	"github.com/qiangyt/jog/util"
)

// FieldMapT ...
type FieldMapT struct {
	Unknown           UnknownFields
	standardsOriginal map[string]Field
	Standards         map[string]Field
}

// FieldMap ...
type FieldMap = *FieldMapT

// Reset ...
func (i FieldMap) Reset() {
	i.Unknown = &UnknownFieldsT{}
	i.Unknown.Reset()

	i.standardsOriginal = make(map[string]Field)
	i.Standards = make(map[string]Field)
}

// UnmarshalYAML ...
func (i FieldMap) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return UnmarshalYAML(i, unmarshal)
}

// MarshalYAML ...
func (i FieldMap) MarshalYAML() (interface{}, error) {
	return MarshalYAML(i)
}

// Init ...
func (i FieldMap) Init(cfg Configuration) {

}

// FromMap ...
func (i FieldMap) FromMap(m map[string]interface{}) error {
	unknownV := util.ExtractFromMap(m, "unknown")
	if unknownV != nil {
		if err := util.UnmashalYAMLAgain(unknownV, &i.Unknown); err != nil {
			return err
		}
	}

	for k, v := range m {
		//if i.config.HasFieldInPattern(k) == false {

		//}

		var f Field
		if err := util.UnmashalYAMLAgain(v, &f); err != nil {
			return err
		}

		f.Name = k
		i.standardsOriginal[k] = f
		i.Standards[k] = f

		if !f.CaseSensitive {
			lk := strings.ToLower(k)
			if lk != k {
				if old, alreadyHas := i.Standards[lk]; old != f && alreadyHas {
					return fmt.Errorf("duplicated field name: %s", lk)
				}
				i.Standards[lk] = f
			}
		}

		var aliases map[string]bool
		if !f.CaseSensitive {
			aliases = f.Alias.Values
		} else {
			aliases = f.Alias.LowercasedValues
		}
		for aliasName := range aliases {
			if old, alreadyHas := i.Standards[aliasName]; old != f && alreadyHas {
				return fmt.Errorf("duplicated field alias name: %s", aliasName)
			}
			i.Standards[aliasName] = f
		}

		delete(m, k)
	}

	return nil
}

// ToMap ...
func (i FieldMap) ToMap() map[string]interface{} {
	r := make(map[string]interface{})
	r["unknown"] = i.Unknown.ToMap()

	for k, v := range i.standardsOriginal {
		r[k] = v.ToMap()
	}
	return r
}
