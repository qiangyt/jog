package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/pkg/errors"
)

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
	if FileStat(path, false) == nil {
		return false
	}
	return true
}

// MkdirAll ...
func MkdirAll(path string) {
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		panic(errors.Wrapf(err, "failed to create directory: %s", path))
	}
}

// RemoveFile ...
func RemoveFile(path string) {
	if err := os.Remove(path); err != nil {
		panic(errors.Wrapf(err, "failed to delete file: %s", path))
	}
}

// RemoveAllFiles ...
func RemoveAllFiles(path string) {
	if err := os.RemoveAll(path); err != nil {
		panic(errors.Wrapf(err, "failed to delete all files: %s", path))
	}
}

// ReadDir ...
func ReadDir(dir string) []os.FileInfo {
	r, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(errors.Wrapf(err, "failed to read directory: %s", dir))
	}
	return r
}

// ReadFile ...
func ReadFile(path string) []byte {
	r, err := ioutil.ReadFile(path)
	if err != nil {
		panic(errors.Wrap(err, "failed to read file: "+path))
	}
	return r
}

// NowMilliseconds ...
func NowMilliseconds() int64 {
	return time.Now().UnixNano() / 1e6
}

// ToJSON ...
func ToJSON(obj interface{}) string {
	r, err := json.Marshal(obj)
	if err != nil {
		panic(errors.Wrapf(err, "failed to json marshal object: %v (type=%v)", obj, reflect.TypeOf(obj)))
	}
	return string(r)
}

// MapToObject ...
func MapToObject(m map[interface{}]interface{}) map[string]interface{} {
	r := make(map[string]interface{})
	for k, v := range m {
		r[k.(string)] = v
	}
	return r
}

/*
// MapArrayToObjectArray ...
func MapArrayToObjectArray(mapArray []interface{}) []map[string]interface{} {
	r := make([]map[string]interface{})
	for _, m := range mapArray {
		obj := m.(map[interface{}]interface{})
		r = append(r, MapToObject(obj))
	}
	return r
}*/

// NormalizeObjectedMap ...
func NormalizeObjectedMap(obj interface{}) map[string]interface{} {
	r := obj.(map[interface{}]interface{})
	return MapToObject(r)
}

// NormalizeStringArray ...
func NormalizeStringArray(objArray []interface{}) []string {
	r := make([]string, 0)
	for _, obj := range objArray {
		r = append(r, obj.(string))
	}
	return r
}

// ExtractFileName ...
func ExtractFileName(filePath string) string {
	indexOfSlash := strings.LastIndex(filePath, "/")
	if indexOfSlash < 0 {
		return filePath
	}
	return filePath[indexOfSlash+1:]
}

// ExtractFileExt ...
func ExtractFileExt(fileName string) string {
	indexOfDot := strings.LastIndex(fileName, ".")
	if indexOfDot < 0 {
		return ""
	}
	return fileName[indexOfDot+1:]
}

// ExtractFileTitle ...
func ExtractFileTitle(filePath string) string {
	fileName := ExtractFileName(filePath)
	fileExt := ExtractFileExt(fileName)
	return fileName[:(len(fileName) - len(fileExt) - 1)]
}

// StringEndsWith ...
func StringEndsWith(s string, endsWith string) bool {
	pos := strings.LastIndex(s, endsWith)
	if pos < 0 {
		return false
	}
	return pos == len(s)-len(endsWith)
}

// StringStartsWith ...
func StringStartsWith(s string, startsWith string) bool {
	pos := strings.Index(s, startsWith)
	return pos == 0
}
