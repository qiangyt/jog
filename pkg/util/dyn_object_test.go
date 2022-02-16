package _util

import (
	"errors"
	"testing"

	. "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func Test_DynObject4YAML_happy(t *testing.T) {
	assert := require.New(t)

	ctrl := NewController(t)
	defer ctrl.Finish()

	mock := NewMockDynObject(ctrl)

	m := map[string]interface{}{"a": "b"}
	InOrder(
		mock.EXPECT().Reset().Times(1),
		mock.EXPECT().FromMap(m).Times(1),
	)

	err := DynObject4YAML(mock, func(output interface{}) error {
		(*output.(*map[string]interface{}))["a"] = "b"
		return nil
	})

	assert.NoError(err)
}

func Test_DynObject4YAML_error_on_unmarshal(t *testing.T) {
	assert := require.New(t)

	ctrl := NewController(t)
	defer ctrl.Finish()

	mock := NewMockDynObject(ctrl)

	mock.EXPECT().Reset().Times(0)
	mock.EXPECT().FromMap(Any()).Times(0)

	err := DynObject4YAML(mock, func(output interface{}) error {
		return errors.New("expected")
	})

	assert.Error(err)
	assert.Equal("expected", err.Error())
}

func Test_DynObject4YAML_error_on_FromMap(t *testing.T) {
	assert := require.New(t)

	ctrl := NewController(t)
	defer ctrl.Finish()

	mock := NewMockDynObject(ctrl)

	InOrder(
		mock.EXPECT().Reset(),
		mock.EXPECT().FromMap(Any()).Return(errors.New("expected")),
		mock.EXPECT().Reset(),
	)

	err := DynObject4YAML(mock, func(output interface{}) error {
		return nil
	})

	assert.Error(err)
	assert.Equal("expected", err.Error())
}

func Test_DynObject2YAML(t *testing.T) {
	assert := require.New(t)

	ctrl := NewController(t)
	defer ctrl.Finish()

	mock := NewMockDynObject(ctrl)

	mock.EXPECT().ToMap().Times(1).Return(map[string]interface{}{"a": "b"})

	actual, err := DynObject2YAML(mock)

	assert.NoError(err)

	assert.Equal(1, len(actual))
	assert.Equal("b", actual["a"])
}
