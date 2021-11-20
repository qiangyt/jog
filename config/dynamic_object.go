package config

// DynamicObject ...
type DynamicObject interface {
	Reset()
	FromMap(m map[string]interface{}) error
	ToMap() map[string]interface{}
	Init(cfg Configuration)
}

// UnmarshalYAML ...
func UnmarshalYAML(i DynamicObject, unmarshal func(interface{}) error) error {
	m := make(map[string]interface{})
	err := unmarshal(&m)
	if err != nil {
		return err
	}

	i.Reset()
	err = i.FromMap(m)
	if err != nil {
		i.Reset()
	}

	//TODO: check to ensure len(m) == 0, that is, no unknown keys left

	return err
}

// MarshalYAML ...
func MarshalYAML(i DynamicObject) (interface{}, error) {
	return i.ToMap(), nil
}
