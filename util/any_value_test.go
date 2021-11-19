package util

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/stretchr/testify/require"
)

func Test_AnyValueFromRaw_nil(t *testing.T) {
	assert := require.New(t)

	target := AnyValueFromRaw(1, nil, nil)

	assert.Nil(target.Raw)
	assert.True(target.Text == "")
	assert.Equal(1, target.LineNo)
}

func Test_Test_AnyValueFromRaw_Map(t *testing.T) {
	assert := require.New(t)

	raw := map[string]string{"k": "v"}
	target := AnyValueFromRaw(0, raw, nil)

	assert.Equal(raw, target.Raw)
	assert.Equal("{\n  \"k\": \"v\"\n}", target.Text)
	assert.Equal(0, target.LineNo)
}

func Test_Test_AnyValueFromRaw_Array(t *testing.T) {
	assert := require.New(t)

	raw := []string{"k", "v"}
	target := AnyValueFromRaw(-1, raw, nil)

	assert.Equal(raw, target.Raw)
	assert.Equal("[\n  \"k\",\n  \"v\"\n]", target.Text)
	assert.Equal(-1, target.LineNo)

	patches := gomonkey.ApplyFunc(json.MarshalIndent, func(_ interface{}, _ string, _ string) ([]byte, error) {
		return nil, errors.New("")
	})
	defer patches.Reset()

	target = AnyValueFromRaw(-1, raw, nil)
	//TODO: what to assert?
}

func Test_AnyValueFromRaw_Slice(t *testing.T) {
	assert := require.New(t)

	raw := []int{1, 2}[1:]
	target := AnyValueFromRaw(-1, raw, nil)

	assert.Equal(raw, target.Raw)
	assert.Equal("[\n  2\n]", target.Text)
	assert.Equal(-1, target.LineNo)
}

func Test_AnyValueFromRaw_jsonText(t *testing.T) {
	assert := require.New(t)

	raw := "[3]"
	target := AnyValueFromRaw(-1, raw, nil)

	assert.Equal(raw, target.Raw)
	assert.Equal("[\n  3\n]", target.Text)
	assert.Equal(-1, target.LineNo)
}

func Test_AnyValueFromRaw_notJsonText(t *testing.T) {
	assert := require.New(t)

	raw := "[3"
	target := AnyValueFromRaw(-1, raw, nil)

	assert.Equal(raw, target.Raw)
	assert.Equal("[3", target.Text)
	assert.Equal(-1, target.LineNo)
}

func Test_AnyValueFromRaw_invalidKind(t *testing.T) {
	assert := require.New(t)

	raw := struct{ X string }{X: "x"}
	target := AnyValueFromRaw(-1, raw, nil)

	assert.Equal(raw, target.Raw)
	assert.Equal("", target.Text)
	assert.Equal(-1, target.LineNo)
}

func Test_AnyValueFromRaw_replace(t *testing.T) {
	assert := require.New(t)

	raw := "a"
	target := AnyValueFromRaw(0, raw, map[string]string{"a": "a-new"})

	assert.Equal(raw, target.Raw)
	assert.Equal("a-new", target.String())
}
