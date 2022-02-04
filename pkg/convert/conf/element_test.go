package conf

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Enum2_FromMap_ToMap_happy(t *testing.T) {
	assert := require.New(t)

	target := &EnumT{}

	expected := map[string]interface{}{
		"alias": "k1, k2",
		"color": "Red",
	}
	err := target.FromMap(expected)

	assert.NoError(err)
	assert.True(target.Alias.Contains("k1", false))
	assert.True(target.Alias.Contains("k2", false))
	assert.Equal("Red", target.Color.String())

	actual := target.ToMap()

	assert.True("k1, k2" == actual["alias"] || "k2, k1" == actual["alias"])
	assert.Equal("Red", actual["color"])
}
