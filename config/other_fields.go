package config

// OtherFieldsT ...
type OtherFieldsT struct {
	IncludeHiddenFields bool
	Name                Element
	Separator           SeparatorField
	Value               Element
}

// OtherFields ...
type OtherFields = *OtherFieldsT

// Reset ...
func (i OtherFields) Reset() {
	i.IncludeHiddenFields = true

	i.Name = &ElementT{}
	i.Name.Reset()

	i.Separator = &SeparatorFieldT{}
	i.Separator.Reset()

	i.Value = &ElementT{}
	i.Value.Reset()
}

// ToMap ...
func (i OtherFields) ToMap() map[string]interface{} {
	r := make(map[string]interface{})
	r["include-hidden-fields"] = i.IncludeHiddenFields
	r["name"] = i.Name.ToMap()
	r["separator"] = i.Separator.ToMap()
	r["value"] = i.Value.ToMap()
	return r
}
