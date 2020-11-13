package util

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/gookit/goutil/strutil"
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
	if err := os.Remove(path); err != nil {
		panic(errors.Wrapf(err, "failed to delete file: %s", path))
	}
}

// ReadFile ...
func ReadFile(path string) []byte {
	r, err := ioutil.ReadFile(path)
	if err != nil {
		panic(errors.Wrap(err, ""))
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

// MkdirAll ...
func MkdirAll(path string) {
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		panic(errors.Wrapf(err, "failed to create directory: %s", path))
	}
}
