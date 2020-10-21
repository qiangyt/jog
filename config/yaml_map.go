package config

// YAMLMap ...
type YAMLMap interface {
	Reset()
	FromMap(m map[string]interface{}) error
	ToMap() map[string]interface{}
}

// UnmarshalYAML ...
func UnmarshalYAML(i YAMLMap, unmarshal func(interface{}) error) error {
	m := make(map[string]interface{})
	err := unmarshal(&m)
	if err != nil {
		return err
	}

	i.Reset()
	return i.FromMap(m)
}

// MarshalYAML ...
func MarshalYAML(i YAMLMap) (interface{}, error) {
	return i.ToMap(), nil
}
