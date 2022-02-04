package conf

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_CompressPrefixAction_String(t *testing.T) {
	assert := require.New(t)

	assert.Equal("remove-non-first-letter", CompressPrefixAction_RemoveNonFirstLetter.String())
	assert.Equal("remove", CompressPrefixAction_Remove.String())
	assert.Equal("remove-non-first-letter", CompressPrefixAction_Default.String())
	assert.Equal("", (CompressPrefixAction_Default + 99).String())
}

func Test_ParseCompressPrefixAction_(t *testing.T) {
	assert := require.New(t)

	assert.Equal(CompressPrefixAction_RemoveNonFirstLetter, ParseCompressPrefixAction("remove-non-first-letter"))
	assert.Equal(CompressPrefixAction_Remove, ParseCompressPrefixAction("remove"))

	assert.Panics(func() { ParseCompressPrefixAction("wrong") })
}

func Test_CompressPrefix_UnmarshalYAML(t *testing.T) {
	assert := require.New(t)

	called := 0

	target := &CompressPrefixT{}
	err := target.UnmarshalYAML(func(input interface{}) error {
		called += 1
		return nil
	})

	assert.NoError(err)
	assert.Equal(1, called)
}

func Test_CompressPrefix_FromMap_ToMap_happy(t *testing.T) {
	assert := require.New(t)

	target := &CompressPrefixT{}
	target.Reset()
	err := target.FromMap(map[string]interface{}{})

	assert.NoError(err)
	assert.False(target.Enabled)
	assert.NotNil(target.Separators.IsEmpty())
	assert.NotNil(target.WhiteList.IsEmpty())
	assert.Equal(CompressPrefixAction_Default, target.Action)

	actual := target.ToMap()

	assert.False(actual["enabled"].(bool))
	assert.Equal("", actual["separators"])
	assert.Equal("", actual["white-list"])
	assert.Equal("remove-non-first-letter", actual["action"])

	target.Reset()
	err = target.FromMap(map[string]interface{}{
		"enabled":    true,
		"separators": ".",
		"white-list": "com.",
		"action":     "remove",
	})

	assert.NoError(err)
	assert.True(target.Enabled)
	assert.True(target.Separators.Contains("."))
	assert.True(target.WhiteList.Contains("com."))
	assert.Equal(CompressPrefixAction_Remove, target.Action)

	actual = target.ToMap()

	assert.True(actual["enabled"].(bool))
	assert.Equal(".", actual["separators"])
	assert.Equal("com.", actual["white-list"])
	assert.Equal("remove", actual["action"])
}

func Test_CompressPrefix_RemoveNonFirstLetter(t *testing.T) {
	assert := require.New(t)

	target := &CompressPrefixT{}
	target.Reset()
	target.Enabled = true
	target.Separators.Parse(".")
	target.WhiteList.Parse("com.")
	target.Action = CompressPrefixAction_RemoveNonFirstLetter

	// white-list-ed
	assert.Equal("com.example", target.Compress("com.example"))

	// no separator
	assert.Equal("comexample", target.Compress("comexample"))

	// has separator
	assert.Equal("o.example", target.Compress("org.example"))
	assert.Equal("o.e.app", target.Compress("org.example.app"))

	// cached
	assert.Equal("o.e.app", target.Compress("org.example.app"))
}

func Test_CompressPrefix_Remove(t *testing.T) {
	assert := require.New(t)

	target := &CompressPrefixT{}
	target.Reset()
	target.Enabled = true
	target.Separators.Parse(".")
	target.WhiteList.Parse("com.")
	target.Action = CompressPrefixAction_Remove

	// white-list-ed
	assert.Equal("com.example", target.Compress("com.example"))

	// no separator
	assert.Equal("comexample", target.Compress("comexample"))

	// has separator
	assert.Equal("example", target.Compress("org.example"))
	assert.Equal("app", target.Compress("org.example.app"))

	// cached
	assert.Equal("app", target.Compress("org.example.app"))
}

func Test_CompressPrefix_Other(t *testing.T) {
	assert := require.New(t)

	target := &CompressPrefixT{}
	target.Reset()
	target.Enabled = true
	target.Separators.Parse(".")
	target.WhiteList.Parse("com.")
	target.Action = -1

	assert.Equal("org.example.app", target.Compress("org.example.app"))
}

func Test_CompressPrefix_detectSeparator(t *testing.T) {
	assert := require.New(t)

	target := &CompressPrefixT{}
	target.Reset()
	target.Separators.CaseSensitive = false

	target.Separators.Parse("a")
	separator, separated := target.detectSeparator("1A2")
	assert.Equal("A", separator)
	assert.Equal([]string{"1", "2"}, separated)

	target.Separators.Parse("A")
	separator, separated = target.detectSeparator("1a2")
	assert.Equal("a", separator)
	assert.Equal([]string{"1", "2"}, separated)
}
