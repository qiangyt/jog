package _util

import (
	"fmt"
	"strings"

	"github.com/gookit/color"
	"github.com/gookit/goutil/strutil"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// UnmashalYAMLAgain ...
func UnmashalYAMLAgain(in interface{}, out interface{}) error {
	yml, err := yaml.Marshal(in)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(yml, out)
}

// ToBool ...
func ToBool(v interface{}) bool {
	switch v.(type) {
	case bool:
		return v.(bool)
	default:
		return strutil.MustBool(strutil.MustString(v))
	}
}

// ExtractFromMap ...
func ExtractFromMap(m map[string]interface{}, key string) interface{} {
	r, has := m[key]
	if !has {
		return nil
	}
	delete(m, key)
	return r
}

// ExtractStringSliceFromMap ...
func ExtractStringSliceFromMap(m map[string]interface{}, key string) ([]string, error) {
	v, has := m[key]
	if !has || v == nil {
		return []string{}, nil
	}

	r, err := MustStringSlice(v)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse %s: %v", key, v)
	}

	delete(m, key)
	return r, nil
}

// MustStringSlice ...
func MustStringSlice(raw interface{}) ([]string, error) {
	switch raw.(type) {
	case []string:
		return raw.([]string), nil
	case []interface{}:
		{
			r := []string{}
			for _, v := range raw.([]interface{}) {
				r = append(r, v.(string))
			}
			return r, nil
		}
	default:
		return nil, fmt.Errorf("not a string array: %v", raw)
	}
}

func ParseConfigExpression(expr string) (string, string, error) {
	arr := strings.Split(expr, "=")
	if len(arr) != 2 {
		return "", "", fmt.Errorf("invalid config item expression: <%s>", expr)
	}
	return arr[0], arr[1], nil
}

func PrintErrorHint(format string, a ...interface{}) {
	color.Red.Printf(format+". Please check above example\n", a...)
}

// abs function that works for int, Math.Abs only accepts float64
func Abs(value int) int {
	if value < 0 {
		value = value * -1
	}
	return value
}
