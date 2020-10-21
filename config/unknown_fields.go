package config

// UnknownFieldsT ...
type UnknownFieldsT struct {
	Name      Element
	Separator SeparatorField
	Value     Element
}

// UnknownFields ...
type UnknownFields = *UnknownFieldsT

// Reset ...
func (i UnknownFields) Reset() {
	i.Name = &ElementT{}
	i.Name.Reset()

	i.Separator = &SeparatorFieldT{}
	i.Separator.Reset()

	i.Value = &ElementT{}
	i.Value.Reset()
}

// ToMap ...
func (i UnknownFields) ToMap() map[string]interface{} {
	r := make(map[string]interface{})
	r["name"] = i.Name.ToMap()
	r["separator"] = i.Separator.ToMap()
	r["value"] = i.Value.ToMap()
	return r
}
