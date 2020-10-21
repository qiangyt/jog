package config

// Dynamic ...
type Dynamic interface {
	Reset()
	FromMap(m map[string]interface{}) error
	ToMap() map[string]interface{}
	Init(cfg Configuration)
}

// UnmarshalYAML ...
func UnmarshalYAML(i Dynamic, unmarshal func(interface{}) error) error {
	m := make(map[string]interface{})
	err := unmarshal(&m)
	if err != nil {
		return err
	}

	i.Reset()
	return i.FromMap(m)
}

// MarshalYAML ...
func MarshalYAML(i Dynamic) (interface{}, error) {
	return i.ToMap(), nil
}
