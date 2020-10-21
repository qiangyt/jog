package config

import (
	"fmt"
	"strings"

	"github.com/qiangyt/jog/util"
)

// FieldMapT ...
type FieldMapT struct {
	Others       OtherFields
	Standards    map[string]Field
	StandardsMap map[string]Field
}

// FieldMap ...
type FieldMap = *FieldMapT

// Reset ...
func (i FieldMap) Reset() {
	i.Others = &OtherFieldsT{}
	i.Others.Reset()

	i.Standards = make(map[string]Field)
	i.StandardsMap = make(map[string]Field)
}

// UnmarshalYAML ...
func (i FieldMap) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return UnmarshalYAML(i, unmarshal)
}

// MarshalYAML ...
func (i FieldMap) MarshalYAML() (interface{}, error) {
	return MarshalYAML(i)
}

// FromMap ...
func (i FieldMap) FromMap(m map[string]interface{}) error {
	othersV := util.ExtractFromMap(m, "others")
	if othersV != nil {
		if err := util.UnmashalYAMLAgain(othersV, &i.Others); err != nil {
			return err
		}
	}

	for k, v := range m {
		var f Field
		if err := util.UnmashalYAMLAgain(v, &f); err != nil {
			return err
		}

		f.Name = k
		i.Standards[k] = f
		i.StandardsMap[k] = f

		if !f.CaseSensitive {
			lk := strings.ToLower(k)
			if lk != k {
				if old, alreadyHas := i.StandardsMap[lk]; old != f && alreadyHas {
					return fmt.Errorf("duplicated field name: %s", lk)
				}
				i.StandardsMap[lk] = f
			}
		}

		var aliases map[string]bool
		if !f.CaseSensitive {
			aliases = f.Alias.Values
		} else {
			aliases = f.Alias.LowercasedValues
		}
		for aliasName := range aliases {
			if old, alreadyHas := i.StandardsMap[aliasName]; old != f && alreadyHas {
				return fmt.Errorf("duplicated field alias name: %s", aliasName)
			}
			i.StandardsMap[aliasName] = f
		}

		delete(m, k)
	}

	return nil
}

// ToMap ...
func (i FieldMap) ToMap() map[string]interface{} {
	r := make(map[string]interface{})
	r["others"] = i.Others.ToMap()

	for k, v := range i.Standards {
		r[k] = v.ToMap()
	}
	return r
}
