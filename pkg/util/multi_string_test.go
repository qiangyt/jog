package util

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

func Test_MultiString_Set_and_Reset(t *testing.T) {
	assert := require.New(t)

	target := &MultiStringT{}

	target.Set("A")
	assert.Equal("A", target.Text)
	assert.Equal(1, len(target.LowercasedValues))
	assert.True(target.LowercasedValues["a"])
	assert.Equal(1, len(target.Values))
	assert.True(target.Values["A"])

	text := " A, \nb, \tC \r"
	target.Set(text)
	assert.Equal(text, target.Text)
	assert.Equal(3, len(target.LowercasedValues))
	assert.True(target.LowercasedValues["a"])
	assert.True(target.LowercasedValues["b"])
	assert.True(target.LowercasedValues["c"])
	assert.Equal(3, len(target.Values))
	assert.True(target.Values["A"])
	assert.True(target.Values["b"])
	assert.True(target.Values["C"])

	target.Reset()
	assert.Equal("", target.Text)
	assert.Equal(0, len(target.LowercasedValues))
	assert.Equal(0, len(target.Values))
}

func Test_MultiString_Containes(t *testing.T) {
	assert := require.New(t)

	target := &MultiStringT{}

	target.Set("A")

	assert.False(target.Contains("a", true))
	assert.True(target.Contains("A", true))

	assert.True(target.Contains("a", false))
	assert.True(target.Contains("A", false))

	assert.False(target.Contains("b", false))
	assert.False(target.Contains("B", true))
}

func Test_MultiString_MarshalYAML(t *testing.T) {
	assert := require.New(t)

	target := &MultiStringT{}

	target.Set("A")
	ymlBytes, err := yaml.Marshal(target)
	assert.Equal("A\n", string(ymlBytes))
	assert.NoError(err)

	target.Set(" A, b, C")
	ymlBytes, err = yaml.Marshal(target)
	assert.Equal("' A, b, C'\n", string(ymlBytes)) //TODO: is this a bug?
	assert.NoError(err)
}

func Test_MultiString_UnmarshalYAML(t *testing.T) {
	assert := require.New(t)

	target := &MultiStringT{}

	yml := "12, 34"
	err := yaml.Unmarshal([]byte(yml), &target)
	assert.NoError(err)
	assert.Equal(2, len(target.Values))
	assert.True(target.Contains("12", true))
	assert.True(target.Contains("34", true))

	err = target.UnmarshalYAML(func(_ interface{}) error {
		return errors.New("")
	})
	assert.Error(err)

	// no changes
	assert.Equal(2, len(target.Values))
	assert.True(target.Contains("12", true))
	assert.True(target.Contains("34", true))

}
