package util

// DynObject ...
type DynObject interface {
	Reset()
	FromMap(m map[string]interface{}) error
	ToMap() map[string]interface{}
}

// DynObject4YAML ...
func DynObject4YAML(i DynObject, unmarshal func(interface{}) error) error {
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

// DynObject2YAML ...
func DynObject2YAML(i DynObject) (map[string]interface{}, error) {
	return i.ToMap(), nil
}
