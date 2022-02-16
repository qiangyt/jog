package _util

import (
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/gookit/color"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

func Test_ColorsFromLabel_happy(t *testing.T) {
	require := require.New(t)

	r, err := ColorsFromLabel("Red")
	require.NoError(err)

	require.Equal(color.Red.Code(), r.Code())

	r, err = ColorsFromLabel("Red, Green")
	require.NoError(err)

	require.Equal(2, len(r))
	require.Equal(color.New(color.Red, color.Green).Code(), r.Code())
}

func Test_ColorsFromLabel_unknownLabel(t *testing.T) {
	_, err := ColorsFromLabel("Red,Xyz")
	require.Error(t, err)
}

func Test_Color_Reset(t *testing.T) {
	require := require.New(t)

	r := &ColorT{}
	r.Reset()

	require.Equal("FgDefault", r.label)
	require.Equal(color.FgDefault, r.style[0])
}

func Test_Colors_Set_happy(t *testing.T) {
	r := &ColorT{}

	r.Set("Green")
	require.Equal(t, color.Green.Code(), r.style.Code())
}

func Test_Colors_Set_fail(t *testing.T) {
	r := &ColorT{}
	require.Panics(t, func() {
		r.Set("wrong")
	})
}

func Test_Colors_UnmarshalYAML_happy(t *testing.T) {
	require := require.New(t)

	r := &ColorT{}

	err := yaml.Unmarshal([]byte("Blue"), &r)
	require.NoError(err)

	require.Equal(color.Blue.Code(), r.style.Code())
}

func Test_Colors_UnmarshalYAML_failed_due_to_invalid_color(t *testing.T) {
	require := require.New(t)

	r := &ColorT{}

	err := yaml.Unmarshal([]byte("RedWrong"), &r)
	if err == nil {
		require.FailNow("expect unmarchal failure but nothing happened")
	}

	require.Equal(color.FgDefault.Code(), r.style.Code())
}

func Test_Colors_UnmarshalYAML_failed_due_to_invalid_yaml(t *testing.T) {
	require := require.New(t)

	r := &ColorT{}

	err := yaml.Unmarshal([]byte("wrong:"), &r)
	if err == nil {
		require.FailNow("expect unmarchal failure but nothing happened")
	}

	require.Equal(color.FgDefault.Code(), r.style.Code())
}

func Test_Colors_MarshalYAML_happy(t *testing.T) {
	require := require.New(t)

	r := &ColorT{}
	r.Set("Red,Green,Blue")
	yamlText, err := r.MarshalYAML()

	require.NoError(err)
	require.Equal("Red,Green,Blue", yamlText)
}

func Test_Colors_Sprint(t *testing.T) {
	assert := require.New(t)

	r := &ColorT{}

	patches := gomonkey.ApplyMethod(reflect.TypeOf(r.style), "Sprint", func(s color.Style, a ...interface{}) string {
		assert.Equal("x", a[0])
		return "called"
	})
	defer patches.Reset()

	r.Set("Red")
	assert.Equal("called", r.Sprint("x"))
}

func Test_Colors_Sprintf(t *testing.T) {
	assert := require.New(t)

	r := &ColorT{}

	patches := gomonkey.ApplyMethod(reflect.TypeOf(r.style), "Sprintf", func(s color.Style, format string, a ...interface{}) string {
		assert.Equal("%s", format)
		assert.Equal("x", a[0])
		return "called"
	})
	defer patches.Reset()

	r.Set("Red")
	assert.Equal("called", r.Sprintf("%s", "x"))
}
