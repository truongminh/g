package s2

import (
	"github.com/golang/geo/s2"
	"sort"
)

func (pq s2CellQueue) ascendRange(min s2.CellID, max s2.CellID, f func(IEntry)) {
	var index = -1
	sort.Search(pq.Len(), func(i int) bool {
		if pq[i].S2CellID() > min {
			// the candidate
			index = i
			return true
		}
		return false
	})

	if index == -1 {
		return
	}

	for i := index; i < len(pq); i++ {
		if pq[i].S2CellID() < max {
			f(pq[i])
		} else {
			break
		}
	}
}
