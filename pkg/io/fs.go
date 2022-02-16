package _io

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
)

// ExeDirectory ...
func ExeDirectory() string {
	return exeDirectory(os.Args[0])
}

// ExeDirectory_ ...
func exeDirectory(exePath string) string {
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

// ReadTextFile ...
func ReadTextFile(path string) string {
	return string(ReadFile(path))
}

// ReadAll ...
func ReadAll(reader io.Reader) []byte {
	r, err := ioutil.ReadAll(reader)
	if err != nil {
		panic(errors.Wrapf(err, "failed to read from Reader: %v", reader))
	}
	return r
}

// ReadAllText ...
func ReadAllText(reader io.Reader) string {
	return string(ReadAll(reader))
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

// ExpandHomePath ...
func ExpandHomePath(path string) string {
	var r string
	var err error

	if r, err = homedir.Expand(path); err != nil {
		panic(errors.Wrapf(err, "failed to expand path: %s", path))
	}
	return r
}

// MkdirAll ...
func MkdirAll(path string) {
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		panic(errors.Wrapf(err, "failed to create directory: %s", path))
	}
}

const (
	jogHomeDir = "~/.jog"
)

// JogHomeDir ...
func JogHomeDir(expand bool, children ...string) string {
	var r string

	if !expand {
		r = jogHomeDir
	} else {
		r = ExpandHomePath(jogHomeDir)
		MkdirAll(r)
	}

	return filepath.Join(r, filepath.Join(children...))
}
