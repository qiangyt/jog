package convert

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func Test_StringSet_Parse_string(t *testing.T) {
	assert := require.New(t)

	target := &StringSetT{}

	target.CaseSensitive = false
	target.Parse("a, B")
	s := target.String()
	assert.True("a, B" == s || "B, a" == s)

	assert.Equal(2, len(target.ValueMap))
	assert.True(target.ValueMap["a"])
	assert.True(target.ValueMap["B"])

	assert.Equal(2, len(target.LowercasedValueMap))
	assert.True(target.LowercasedValueMap["a"])
	assert.True(target.LowercasedValueMap["b"])

	assert.Equal(2, len(target.UppercasedValueMap))
	assert.True(target.UppercasedValueMap["A"])
	assert.True(target.UppercasedValueMap["B"])

	target.CaseSensitive = true
	target.Parse("c, D")
	s = target.String()
	assert.True("c, D" == s || "D, c" == s)

	assert.Equal(2, len(target.ValueMap))
	assert.True(target.ValueMap["c"])
	assert.True(target.ValueMap["D"])
	assert.Equal(0, len(target.LowercasedValueMap))
	assert.Equal(0, len(target.UppercasedValueMap))
}

func Test_StringSet_Parse_string_array(t *testing.T) {
	assert := require.New(t)

	target := &StringSetT{}

	target.CaseSensitive = false
	target.Parse([]string{" A ", " b"})
	s := target.String()
	assert.True("A, b" == s || "b, A" == s)

	assert.Equal(2, len(target.ValueMap))
	assert.True(target.ValueMap["A"])
	assert.True(target.ValueMap["b"])
}

func Test_StringSet_Parse_fail(t *testing.T) {
	assert := require.New(t)

	target := &StringSetT{}
	assert.Panics(func() {
		target.Parse([]int{1, 2})
	})
}

func Test_StringSet_Contains(t *testing.T) {
	assert := require.New(t)

	target := &StringSetT{}

	target.CaseSensitive = false
	target.Parse("X")
	assert.Equal("X", target.String())

	assert.True(target.Contains("X"))
	assert.True(target.Contains("x"))
	assert.False(target.Contains("y"))

	target.CaseSensitive = true
	target.Parse("y")
	assert.Equal("y", target.String())

	assert.True(target.Contains("y"))
	assert.False(target.Contains("Y"))
	assert.False(target.Contains("x"))
}

func Test_StringSet_ContainsPrefixOf(t *testing.T) {
	assert := require.New(t)

	target := &StringSetT{}

	assert.False(target.ContainsPrefixOf(""))

	target.CaseSensitive = false
	target.Parse("X")
	assert.Equal("X", target.String())
	assert.True(target.ContainsPrefixOf("X123"))
	assert.True(target.ContainsPrefixOf("x123"))
	assert.False(target.ContainsPrefixOf("123X"))

	target.CaseSensitive = true
	target.Parse("y123")
	assert.Equal("y123", target.String())
	assert.True(target.ContainsPrefixOf("y123a"))
	assert.False(target.ContainsPrefixOf("Y123a"))
	assert.False(target.ContainsPrefixOf("123YB"))
}

func Test_StringSet_UnmarshalYAML_happy(t *testing.T) {
	assert := require.New(t)

	called := 0

	target := &StringSetT{CaseSensitive: true}
	err := target.UnmarshalYAML(func(input interface{}) error {
		*input.(*interface{}) = "A"
		called += 1
		return nil
	})

	assert.NoError(err)
	assert.Equal(1, called)
	assert.True(target.Contains("A"))
	assert.False(target.Contains("a"))
}

func Test_StringSet_UnmarshalYAML_fail(t *testing.T) {
	assert := require.New(t)

	target := &StringSetT{}
	err := target.UnmarshalYAML(func(input interface{}) error {
		return errors.New("expected")
	})

	assert.Error(err)
	assert.Equal("expected", err.Error())
}

func Test_StringSet_MarshalYAML(t *testing.T) {
	assert := require.New(t)

	target := &StringSetT{}
	target.Parse("x")
	yamlText, err := target.MarshalYAML()

	assert.NoError(err)
	assert.Equal("x", yamlText)
}
