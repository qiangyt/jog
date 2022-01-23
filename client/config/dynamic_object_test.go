package config

import (
	"errors"
	"testing"

	. "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func Test_UnmarshalYAML_happy(t *testing.T) {
	assert := require.New(t)

	ctrl := NewController(t)
	defer ctrl.Finish()

	mock := NewMockDynamicObject(ctrl)

	m := map[string]interface{}{"a": "b"}
	InOrder(
		mock.EXPECT().Reset().Times(1),
		mock.EXPECT().FromMap(m).Times(1),
	)

	err := UnmarshalYAML(mock, func(output interface{}) error {
		(*output.(*map[string]interface{}))["a"] = "b"
		return nil
	})

	assert.NoError(err)
}

func Test_UnmarshalYAML_error_on_unmarshal(t *testing.T) {
	assert := require.New(t)

	ctrl := NewController(t)
	defer ctrl.Finish()

	mock := NewMockDynamicObject(ctrl)

	mock.EXPECT().Reset().Times(0)
	mock.EXPECT().FromMap(Any()).Times(0)

	err := UnmarshalYAML(mock, func(output interface{}) error {
		return errors.New("expected")
	})

	assert.Error(err)
	assert.Equal("expected", err.Error())
}

func Test_UnmarshalYAML_error_on_FromMap(t *testing.T) {
	assert := require.New(t)

	ctrl := NewController(t)
	defer ctrl.Finish()

	mock := NewMockDynamicObject(ctrl)

	InOrder(
		mock.EXPECT().Reset(),
		mock.EXPECT().FromMap(Any()).Return(errors.New("expected")),
		mock.EXPECT().Reset(),
	)

	err := UnmarshalYAML(mock, func(output interface{}) error {
		return nil
	})

	assert.Error(err)
	assert.Equal("expected", err.Error())
}

func Test_MarshalYAML(t *testing.T) {
	assert := require.New(t)

	ctrl := NewController(t)
	defer ctrl.Finish()

	mock := NewMockDynamicObject(ctrl)

	mock.EXPECT().ToMap().Times(1).Return(map[string]interface{}{"a": "b"})

	actual, err := MarshalYAML(mock)

	assert.NoError(err)

	actualMap := actual.(map[string]interface{})
	assert.Equal(1, len(actualMap))
	assert.Equal("b", actualMap["a"])
}
