package util

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
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

func Test_UnmashalYAMLAgain(t *testing.T) {
	assert := require.New(t)

	type Temp struct {
		K string
	}
	in := Temp{K: "v"}
	out := Temp{K: ""}
	assert.NoError(UnmashalYAMLAgain(in, &out))
	assert.Equal(in, out)

	patches := gomonkey.ApplyFunc(yaml.Marshal, func(_ interface{}) ([]byte, error) {
		return nil, errors.New("")
	})
	defer patches.Reset()
	assert.Error(UnmashalYAMLAgain(in, &out))
}

func Test_ToBool(t *testing.T) {
	assert := require.New(t)

	assert.True(ToBool(true))
	assert.False(ToBool(false))

	assert.True(ToBool("true"))
	assert.False(ToBool("false"))
	assert.True(ToBool("True"))
	assert.False(ToBool("False"))
	assert.True(ToBool("TRUE"))
	assert.False(ToBool("FALSE"))

	assert.True(ToBool("yes"))
	assert.False(ToBool("no"))
	assert.True(ToBool("Yes"))
	assert.False(ToBool("No"))
	assert.True(ToBool("YES"))
	assert.False(ToBool("NO"))

	assert.True(ToBool(1))
	assert.False(ToBool(0))
	assert.True(ToBool("1"))
	assert.False(ToBool("0"))

	assert.True(ToBool("on"))
	assert.False(ToBool("off"))
	assert.True(ToBool("On"))
	assert.False(ToBool("Off"))
	assert.True(ToBool("ON"))
	assert.False(ToBool("OFF"))

	assert.False(ToBool(nil))
	assert.False(ToBool(map[string]int{}))
}

func Test_ExtractFromMap(t *testing.T) {
	assert := require.New(t)

	m := map[string]interface{}{"k": "v"}

	assert.Equal("v", ExtractFromMap(m, "k"))

	_, has := m["k"]
	assert.False(has, "should be removed")

	assert.Nil(ExtractFromMap(m, "k"))
}

func Test_ExtractStringSliceFromMap(t *testing.T) {
	assert := require.New(t)

	m := map[string]interface{}{
		"k": []string{"v0", "v1"},
		"p": "not slice",
	}

	v, err := ExtractStringSliceFromMap(m, "k")
	assert.NoError(err)
	assert.Equal(2, len(v))
	assert.Equal("v0", v[0])
	assert.Equal("v1", v[1])
	_, has := m["k"]
	assert.False(has, "should be removed")

	v, err = ExtractStringSliceFromMap(m, "k")
	assert.NoError(err)
	assert.Equal(0, len(v))

	v, err = ExtractStringSliceFromMap(m, "p")
	assert.Error(err)
	assert.Nil(v)
	_, has = m["p"]
	assert.True(has, "should be still there")
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

func Test_MustStringSlice(t *testing.T) {
	assert := require.New(t)

	input := []string{"a"}
	s, err := MustStringSlice(input)
	assert.NoError(err)
	assert.Equal(input, s)

	s, err = MustStringSlice([]interface{}{"A"})
	assert.NoError(err)
	assert.Equal(1, len(s))
	assert.Equal("A", s[0])

	assert.Panics(func() { MustStringSlice([]interface{}{789}) })

	s, err = MustStringSlice([]int{})
	assert.Error(err)
	assert.Nil(s)
}
