package config

// PrefixT ...
type PrefixT struct {
	ElementT
}

// Prefix ...
type Prefix = *PrefixT

// Reset ...
func (i Prefix) Reset() {
	i.ElementT.Reset()

	i.Color.Set("FgBlue")
}

// FromMap ...
func (i Prefix) FromMap(m map[string]interface{}) error {
	return i.ElementT.FromMap(m)
}

// ToMap ...
func (i Prefix) ToMap() map[string]interface{} {
	return i.ElementT.ToMap()
}

// UnmarshalYAML ...
func (i Prefix) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return UnmarshalYAML(i, unmarshal)
}

// MarshalYAML ...
func (i Prefix) MarshalYAML() (interface{}, error) {
	return MarshalYAML(i)
}
