package s2

import (
	"container/heap"
)

// update modifies the priority and value of an Item in the queue.
func (m *Map) Insert(value IEntry) {
	heap.Push(m.q, value)
}

func (m *Map) Update(value IEntry) {
	heap.Fix(m.q, value.S2GetIndex())
}

func (m *Map) Remove(loc IEntry) {
	m.q.Print("remove")
	heap.Remove(m.q, loc.S2GetIndex())
}
