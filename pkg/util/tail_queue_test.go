package _util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_initialization(t *testing.T) {
	assert := require.New(t)

	r := NewTailQueue(3)

	assert.False(r.Count() != 0 || !r.IsEmpty(), "q should be empty when just initialized")
	assert.False(r.IsFull(), "q should be not full when just initialized")
	assert.Nil(r.Head(), "head element should be nil")
	assert.Nil(r.Tail(), "tail element should be nil")
}

func Test_Add_aLittle(t *testing.T) {
	assert := require.New(t)

	r := NewTailQueue(3)

	r.Add("1")
	r.Add("2")
	assert.False(r.Count() != 2 || r.IsEmpty(), "element amount should be 2")
	assert.False(r.IsFull(), "q(3) with 2 elements amount should not be full")
	assert.Equal("1", r.Head(), "head element should be '1'")
	assert.Equal("2", r.Tail(), "tail element should be '2'")

	e1 := r.Kick()
	assert.Equal("1", e1, "the kicked element 1 should be '1'")
	assert.False(r.Count() != 1 || r.IsEmpty(), "after kicked one, element amount should be 1")
	assert.False(r.IsFull(), "after kicked one, q(3) with 2 elements amount should not be full")
	assert.Equal("2", r.Head(), "after kicked one, head element should be '2'")
	assert.False(r.Tail() != r.Head(), "after kicked one, tail element and head element should be same")

	e2 := r.Kick()
	assert.Equal("2", e2, "the kicked element 2 should be '2'")
	assert.False(r.Count() != 0 || !r.IsEmpty(), "after kicked two, element amount should be 0")
	assert.False(r.IsFull(), "after kicked two, q(3) with 2 elements amount should not be full")
	assert.Nil(r.Head(), "after kicked two, head element should be nil")
	assert.Nil(r.Tail(), "after kicked two, tail element should be nil")
}

func Test_Add_Full_then_Kick(t *testing.T) {
	assert := require.New(t)

	r := NewTailQueue(3)

	r.Add("1")
	r.Add("2")
	r.Add("3")

	assert.False(r.Count() != 3 || r.IsEmpty(), "after add 3 element, should not be empty")
	assert.True(r.IsFull(), "after add 3 element, should be full")
	assert.Equal("1", r.Head(), "head element should be '1'")
	assert.Equal("3", r.Tail(), "tail element should be '3'")

	r.Add("4")
	assert.False(r.Count() != 3 || r.IsEmpty(), "after add 4 element, should not be empty")
	assert.True(r.IsFull(), "after add 4 element, should be full")
	assert.Equal("2", r.Head(), "head element should be '2'")
	assert.Equal("4", r.Tail(), "tail element should be '4'")

	r.Add("5")
	r.Add("6")
	r.Add("7")

	assert.False(r.Count() != 3 || r.IsEmpty(), "after add 7 element, should not be empty")
	assert.True(r.IsFull(), "after add 7 element, should be full")
	assert.Equal("5", r.Head(), "head element should be '5'")
	assert.Equal("7", r.Tail(), "tail element should be '7'")

	r.Add("8")
	r.Add("9")
	assert.Equal(3, r.Count(), "element count should be 3")

	assert.Equal("7", r.Kick(), "should be '7'")
	assert.Equal(2, r.Count(), "element count should be 2")
	assert.Equal("8", r.Kick(), "should be '8'")
	assert.Equal("9", r.Kick(), "should be '9'")

	assert.Equal(0, r.Count(), "element count should be 0")
	assert.True(r.IsEmpty(), "q should be empty")

	assert.False(r.Kick() != nil || r.Head() != nil || r.Tail() != nil, "should no element")
}

func Test_Clear(t *testing.T) {
	assert := require.New(t)

	r := NewTailQueue(3)

	r.Add("1")
	r.Add("2")

	assert.False(r.Count() != 2 || r.IsEmpty(), "element amount should be 2")

	r.Clear()

	assert.False(r.Count() != 0 || !r.IsEmpty(), "after clear, should be empty")
	assert.False(r.IsFull(), "after clear, should not be full")
	assert.Nil(r.Head(), "after clear, head element should be nil")
	assert.Nil(r.Tail(), "after clear, tail element should be nil")
}
