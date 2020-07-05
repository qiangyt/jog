package util

// YAMLMap ...
type YAMLMap interface {
	Reset()
	FromMap(m map[string]interface{}) error
	ToMap(m map[string]interface{}) error
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
	r := make(map[string]interface{})
	err := i.ToMap(r)
	return r, err
}
