package ex

import "time"

type IEntry interface {
	ExMTime() int
	ExIndex() int
	ExSetIndex(i int)
}

type ExEntry struct {
	exMTime int
	exIndex int
}

func (e *ExEntry) ExIndex() int {
	return e.exIndex
}

func (e *ExEntry) ExSetIndex(i int) {
	e.exIndex = i
}

func (e *ExEntry) ExMTime() int {
	return e.exMTime
}

func (e *ExEntry) ExUpdate() {
	var now = time.Now().Unix()
	e.exMTime = int(now)
}
