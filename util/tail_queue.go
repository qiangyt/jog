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
func (i TailQueue) Clear() {
	i.frontIndex = 0
	i.rearIndex = 0
	i.count = 0
}

// Count ...
func (i TailQueue) Count() int {
	return i.count
}

// IsEmpty ...
func (i TailQueue) IsEmpty() bool {
	return i.count == 0
}

// IsFull ...
func (i TailQueue) IsFull() bool {
	return i.count == len(i.items)
}

// Add ...
func (i TailQueue) Add(element interface{}) {
	if i.IsFull() {
		i.Kick()
	}

	i.items[i.rearIndex] = element
	i.rearIndex = (i.rearIndex + 1) % len(i.items)
	i.count = i.count + 1
}

// Kick ...
func (i TailQueue) Kick() interface{} {
	if i.IsEmpty() {
		return nil
	}

	r := i.items[i.frontIndex]

	i.frontIndex = (i.frontIndex + 1) % len(i.items)
	i.count = i.count - 1

	return r
}
