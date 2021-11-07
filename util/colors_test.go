package util

import (
	"testing"

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

func Test_Colors_UnmarshalYAML_failed(t *testing.T) {
	require := require.New(t)

	r := &ColorT{}

	err := yaml.Unmarshal([]byte("Nothing"), &r)
	if err == nil {
		require.FailNow("expect unmarchal failure but nothing happened")
	}

	require.Equal(color.FgDefault.Code(), r.style.Code())
}
