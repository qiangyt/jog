package util

// TailQueueT ...
type TailQueueT struct {
	frontIndex int
	rearIndex  int
	count      int
	items      []interface{}
}

// TailQueue ...
type TailQueue = *TailQueueT

// NewTailQueueWithSize ...
func NewTailQueueWithSize(size int) TailQueue {
	r := &TailQueueT{items: make([]interface{}, size)}
	r.Clear()
	return r
}

// Clear ...
func (me TailQueue) Clear() {
	me.frontIndex = 0
	me.rearIndex = 0
	me.count = 0
}

// Count ...
func (me TailQueue) Count() int {
	return me.count
}

// IsEmpty ...
func (me TailQueue) IsEmpty() bool {
	return me.count == 0
}

// IsFull ...
func (me TailQueue) IsFull() bool {
	return me.count == len(me.items)
}

// Add ...
func (me TailQueue) Add(element interface{}) {
	if me.IsFull() {
		me.Kick()
	}

	me.items[me.rearIndex] = element
	me.rearIndex = (me.rearIndex + 1) % len(me.items)
	me.count = me.count + 1
}

// Kick ...
func (me TailQueue) Kick() interface{} {
	if me.IsEmpty() {
		return nil
	}

	r := me.items[me.frontIndex]

	me.frontIndex = (me.frontIndex + 1) % len(me.items)
	me.count = me.count - 1

	return r
}
