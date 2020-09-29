package util

import "testing"

func Test_initialization(t *testing.T) {
	r := NewTailQueue(3)

	if r.Count() != 0 || !r.IsEmpty() {
		t.Error("q should be empty when just initialized")
	}

	if r.IsFull() {
		t.Error("q should be not full when just initialized")
	}

	if r.Head() != nil {
		t.Error("head element should be nil")
	}
	if r.Tail() != nil {
		t.Error("tail element should be nil")
	}
}

func Test_Add_aLittle(t *testing.T) {
	r := NewTailQueue(3)

	r.Add("1")
	r.Add("2")
	if r.Count() != 2 || r.IsEmpty() {
		t.Error("element amount should be 2")
	}
	if r.IsFull() {
		t.Error("q(3) with 2 elements amount should not be full")
	}
	if r.Head() != "1" {
		t.Error("head element should be '1'")
	}
	if r.Tail() != "2" {
		t.Error("tail element should be '2'")
	}

	e1 := r.Kick()
	if e1 != "1" {
		t.Error("the kicked element 1 should be '1'")
	}
	if r.Count() != 1 || r.IsEmpty() {
		t.Error("after kicked one, element amount should be 1")
	}
	if r.IsFull() {
		t.Error("after kicked one, q(3) with 2 elements amount should not be full")
	}
	if r.Head() != "2" {
		t.Error("after kicked one, head element should be '2'")
	}
	if r.Tail() != r.Head() {
		t.Error("after kicked one, tail element and head element should be same")
	}

	e2 := r.Kick()
	if e2 != "2" {
		t.Error("the kicked element 2 should be '2'")
	}
	if r.Count() != 0 || !r.IsEmpty() {
		t.Error("after kicked two, element amount should be 0")
	}
	if r.IsFull() {
		t.Error("after kicked two, q(3) with 2 elements amount should not be full")
	}
	if r.Head() != nil {
		t.Error("after kicked two, head element should be nil")
	}
	if r.Tail() != nil {
		t.Error("after kicked two, tail element should be nil")
	}

}

func Test_Add_Full_then_Kick(t *testing.T) {
	r := NewTailQueue(3)

	r.Add("1")
	r.Add("2")
	r.Add("3")

	if r.Count() != 3 || r.IsEmpty() {
		t.Error("after add 3 element, should not be empty")
	}
	if !r.IsFull() {
		t.Error("after add 3 element, should be full")
	}
	if r.Head() != "1" {
		t.Error("head element should be '1'")
	}
	if r.Tail() != "3" {
		t.Error("tail element should be '3'")
	}

	r.Add("4")
	if r.Count() != 3 || r.IsEmpty() {
		t.Error("after add 4 element, should not be empty")
	}
	if !r.IsFull() {
		t.Error("after add 4 element, should be full")
	}
	if r.Head() != "2" {
		t.Error("head element should be '2'")
	}
	if r.Tail() != "4" {
		t.Error("tail element should be '4'")
	}

	r.Add("5")
	r.Add("6")
	r.Add("7")

	if r.Count() != 3 || r.IsEmpty() {
		t.Error("after add 7 element, should not be empty")
	}
	if !r.IsFull() {
		t.Error("after add 7 element, should be full")
	}
	if r.Head() != "5" {
		t.Error("head element should be '5'")
	}
	if r.Tail() != "7" {
		t.Error("tail element should be '7'")
	}

	r.Add("8")
	r.Add("9")
	if r.Count() != 3 {
		t.Error("element count should be 3")
	}

	if r.Kick() != "7" {
		t.Error("should be '7'")
	}
	if r.Count() != 2 {
		t.Error("element count should be 2")
	}
	if r.Kick() != "8" {
		t.Error("should be '8'")
	}
	if r.Kick() != "9" {
		t.Error("should be '9'")
	}

	if r.Count() != 0 {
		t.Error("element count should be 0")
	}
	if !r.IsEmpty() {
		t.Error("q should be empty")
	}

	if r.Kick() != nil || r.Head() != nil || r.Tail() != nil {
		t.Error("should no element")
	}
}

func Test_Clear(t *testing.T) {
	r := NewTailQueue(3)

	r.Add("1")
	r.Add("2")

	if r.Count() != 2 || r.IsEmpty() {
		t.Error("element amount should be 2")
	}

	r.Clear()

	if r.Count() != 0 || !r.IsEmpty() {
		t.Error("after clear, should be empty")
	}
	if r.IsFull() {
		t.Error("after clear, should not be full")
	}
	if r.Head() != nil {
		t.Error("after clear, head element should be nil")
	}
	if r.Tail() != nil {
		t.Error("after clear, tail element should be nil")
	}

}
