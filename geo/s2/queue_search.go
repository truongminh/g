package s2

import (
	"github.com/golang/geo/s2"
	"sort"
)

func (pq s2CellQueue) ascendRange(min s2.CellID, max s2.CellID, f func(IEntry)) {
	var index = sort.Search(pq.Len(), func(i int) bool {
		return pq[i].S2CellID() > min
	})

	for i := index; i < len(pq); i++ {
		if pq[i].S2CellID() < max {
			f(pq[i])
		} else {
			break
		}
	}
}
