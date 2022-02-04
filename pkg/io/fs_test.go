package io

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func Test_ExeDirectory_happy(t *testing.T) {
	assert := require.New(t)

	assert.Contains(ExeDirectory(), "go-build")

	patches := gomonkey.ApplyFunc(filepath.Abs, func(_ string) (string, error) {
		return "", errors.New("")
	})
	defer patches.Reset()

	assert.Panics(func() { exeDirectory("") })
}

func Test_FileStat(t *testing.T) {
	assert := require.New(t)

	fi := FileStat("/home", true)
	assert.NotNil(fi)
	assert.True(fi.IsDir())
	assert.Equal("home", fi.Name())

	assert.Nil(FileStat("/not_existed", false))

	assert.Panics(func() {
		FileStat("/not_exists", true)
	})

	assert.Panics(func() {
		FileStat("/root/.ssh", false)
	})
}

func Test_FileExists(t *testing.T) {
	assert := require.New(t)

	assert.False(FileExists("/not_exists.txt"))
	assert.True(FileExists("/etc/hosts"))

	assert.Panics(func() {
		FileExists("/home")
	})
}

func Test_DirExists(t *testing.T) {
	assert := require.New(t)

	assert.False(DirExists("/not_exists"))
	assert.True(DirExists("/home"))

	assert.Panics(func() {
		DirExists("/etc/hosts")
	})
}

func Test_RemoveFile(t *testing.T) {
	assert := require.New(t)

	f, _ := ioutil.TempFile("", "jog_*.txt")
	path := f.Name()
	f.Close()

	RemoveFile(path)

	WriteFileIfNotFound(path, []byte("abc"))
	assert.True(FileExists(path))
	RemoveFile(path)
	assert.False(FileExists(path))

	assert.True(FileExists("/etc/hosts"))
	assert.Panics(func() {
		RemoveFile("/etc/hosts")
	})
	assert.True(FileExists("/etc/hosts"))
}

func Test_RemoveDir(t *testing.T) {
	assert := require.New(t)

	path, _ := ioutil.TempDir("", "jog_*")
	assert.True(DirExists(path))
	RemoveDir(path)
	assert.False(DirExists(path))

	assert.True(DirExists("/root"))
	assert.Panics(func() {
		RemoveDir("/root")
	})
	assert.True(DirExists("/root"))

	assert.True(DirExists("/"))
	assert.Panics(func() {
		RemoveDir("/")
	})
	assert.True(DirExists("/"))
}

func Test_ReadFile(t *testing.T) {
	assert := require.New(t)

	f, _ := ioutil.TempFile("", "jog_*.txt")
	path := f.Name()
	f.Close()

	RemoveFile(path)
	assert.Panics(func() {
		ReadFile(path)
	})

	WriteFileIfNotFound(path, []byte("abc"))
	assert.True(FileExists(path))

	content := ReadFile(path)
	assert.Equal("abc", string(content))
}

func Test_WriteFileIfNotFound(t *testing.T) {
	assert := require.New(t)

	f, _ := ioutil.TempFile("", "jog_*.txt")
	path := f.Name()
	f.Close()

	RemoveFile(path)
	WriteFileIfNotFound(path, []byte("old"))
	content := ReadFile(path)
	assert.Equal("old", string(content))

	WriteFileIfNotFound(path, []byte("new"))
	content = ReadFile(path)
	assert.Equal("old", string(content))
}

func Test_ReplaceFile(t *testing.T) {
	assert := require.New(t)

	f, _ := ioutil.TempFile("", "jog_*.txt")
	path := f.Name()
	f.Close()

	RemoveFile(path)
	ReplaceFile(path, []byte("1"))
	content := ReadFile(path)
	assert.Equal("1", string(content))

	ReplaceFile(path, []byte("2"))
	content = ReadFile(path)
	assert.Equal("2", string(content))

	assert.Panics(func() {
		ReplaceFile("/etc/hosts", []byte("bad"))
	})
}

func Test_ExpandHomePath(t *testing.T) {
	assert := require.New(t)

	home := os.Getenv("HOME")
	assert.Equal(home+"/jog.txt", ExpandHomePath("~/jog.txt"))

	assert.Panics(func() {
		ReadFile(ExpandHomePath("~jog.txt"))
	})
}

func Test_MkdirAll(t *testing.T) {
	assert := require.New(t)

	assert.Panics(func() {
		MkdirAll("/root/test")
	})

	path := os.TempDir() + "/jog_test/1/2/3/4/5/6"
	assert.NoDirExists(path)

	MkdirAll(path)
	defer RemoveDir(path)
	assert.DirExists(path)
}
