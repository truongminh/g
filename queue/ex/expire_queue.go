package ex

import (
	"container/heap"
	"time"
)

type ExpireQueue struct {
	q *expireQueue
}

func NewExpireQueue() *ExpireQueue {
	var q = expireQueue([]IEntry{})
	return &ExpireQueue{q: &q}
}

// update modifies the priority and value of an Item in the queue.
func (e *ExpireQueue) Update(value IEntry) {
	heap.Fix(e.q, value.ExIndex())
}

func (e *ExpireQueue) Insert(value IEntry) {
	heap.Push(e.q, value)
}

func (e *ExpireQueue) remove(value IEntry) {
	heap.Remove(e.q, value.ExIndex())
}

func (e *ExpireQueue) RemoveOld(second int) []IEntry {
	var expireAt = int(time.Now().Unix()) - second
	var removed = make([]IEntry, 0)
	// expire item
	for e.q.Len() > 0 {
		var v = heap.Pop(e.q).(IEntry)
		// pop first?
		if v.ExMTime() > expireAt {
			// repair the pop
			heap.Push(e.q, v)
			break
		} else {
			// e.remove(v)
			removed = append(removed, v)
		}
	}
	return removed
}
