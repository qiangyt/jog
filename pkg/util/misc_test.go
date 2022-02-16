package _util

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

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
		"z": nil,
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

	v, err = ExtractStringSliceFromMap(m, "z")
	assert.NoError(err)
	assert.Equal(0, len(v))
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
