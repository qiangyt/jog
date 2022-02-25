package conf

import (
	"fmt"
	"strings"

	_util "github.com/qiangyt/jog/pkg/util"
)

// FieldMapT ...
type FieldMapT struct {
	Others                  OtherFields
	StandardsWithAllAliases map[string]Field
	Standards               map[string]Field
}

// FieldMap ...
type FieldMap = *FieldMapT

// Reset ...
func (i FieldMap) Reset() {
	i.Others = &OtherFieldsT{}
	i.Others.Reset()

	i.Standards = make(map[string]Field)
	i.StandardsWithAllAliases = make(map[string]Field)
}

// UnmarshalYAML ...
func (i FieldMap) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return _util.DynObject4YAML(i, unmarshal)
}

// MarshalYAML ...
func (i FieldMap) MarshalYAML() (interface{}, error) {
	return _util.DynObject2YAML(i)
}

// FromMap ...
func (i FieldMap) FromMap(m map[string]interface{}) error {
	othersV := _util.ExtractFromMap(m, "others")
	if othersV != nil {
		if err := _util.UnmashalYAMLAgain(othersV, &i.Others); err != nil {
			return err
		}
	}

	for k, v := range m {
		//if i.config.HasFieldInPattern(k) == false {

		//}

		var f Field
		if err := _util.UnmashalYAMLAgain(v, &f); err != nil {
			return err
		}

		if k == "logger" {
			k = "logger"
		}
		f.Name = k
		i.StandardsWithAllAliases[k] = f
		i.Standards[k] = f

		if !f.CaseSensitive {
			lk := strings.ToLower(k)
			if lk != k {
				if old, alreadyHas := i.StandardsWithAllAliases[lk]; old != f && alreadyHas {
					return fmt.Errorf("duplicated field name: %s", lk)
				}
				i.StandardsWithAllAliases[lk] = f
			}
		}

		var aliases map[string]bool
		if !f.CaseSensitive {
			aliases = f.Alias.Values
		} else {
			aliases = f.Alias.LowercasedValues
		}
		for aliasName := range aliases {
			if old, alreadyHas := i.StandardsWithAllAliases[aliasName]; old != f && alreadyHas {
				return fmt.Errorf("duplicated field alias name: %s", aliasName)
			}
			i.StandardsWithAllAliases[aliasName] = f
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