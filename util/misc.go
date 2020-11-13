package util

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/gookit/goutil/strutil"
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// ExeDirectory ...
func ExeDirectory() string {
	exePath := os.Args[0]
	r, err := filepath.Abs(filepath.Dir(exePath))
	if err != nil {
		panic(errors.Wrapf(err, "failed to get absolute directory path for "+exePath))
	}
	return r
}

// FileStat ...
func FileStat(path string, ensureExists bool) os.FileInfo {
	r, err := os.Stat(path)
	if err != nil {
		if !os.IsNotExist(err) {
			panic(errors.Wrapf(err, "failed to stat file: %s", path))
		}
		if ensureExists {
			panic(errors.Wrapf(err, "file not exists: %s", path))
		}
		return nil
	}

	return r
}

// FileExists ...
func FileExists(path string) bool {
	fi := FileStat(path, false)
	if fi == nil {
		return false
	}
	if fi.IsDir() {
		panic(fmt.Errorf("expect %s be file, but it is directory", path))
	}
	return true
}

// DirExists ...
func DirExists(path string) bool {
	fi := FileStat(path, false)
	if fi == nil {
		return false
	}
	if !fi.IsDir() {
		panic(fmt.Errorf("expect %s be directory, but it is file", path))
	}
	return true
}

// RemoveFile ...
func RemoveFile(path string) {
	if FileExists(path) {
		if err := os.Remove(path); err != nil {
			panic(errors.Wrapf(err, "failed to delete file: %s", path))
		}
	}
}

// RemoveDir ...
func RemoveDir(path string) {
	if path == "/" || path == "\\" {
		panic(fmt.Errorf("should NOT remove root directory"))
	}
	if err := os.RemoveAll(path); err != nil {
		panic(errors.Wrapf(err, "failed to delete directory: %s", path))
	}
}

// ReadFile ...
func ReadFile(path string) []byte {
	r, err := ioutil.ReadFile(path)
	if err != nil {
		panic(errors.Wrapf(err, "failed to read file: %s", path))
	}
	return r
}

// WriteFileIfNotFound ...
func WriteFileIfNotFound(path string, content []byte) {
	if FileExists(path) {
		return
	}
	ReplaceFile(path, content)
}

// ReplaceFile ...
func ReplaceFile(path string, content []byte) {
	if err := ioutil.WriteFile(path, content, 0640); err != nil {
		panic(errors.Wrapf(err, "failed to write file: %s", path))
	}
}

// ExpandPath ...
func ExpandPath(path string) string {
	var r string
	var err error

	if r, err = homedir.Expand(path); err != nil {
		panic(errors.Wrapf(err, "failed to expand path: %s", path))
	}
	return r
}

// UnmashalYAMLAgain ...
func UnmashalYAMLAgain(in interface{}, out interface{}) error {
	yml, err := yaml.Marshal(in)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(yml, out)
	return err
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
	v := ExtractFromMap(m, key)
	if v == nil {
		return []string{}, nil
	}

	r, err := MustStringSlice(v)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse %s: %v", key, v)
	}
	return r, nil
}

// MkdirAll ...
func MkdirAll(path string) {
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		panic(errors.Wrapf(err, "failed to create directory: %s", path))
	}
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
