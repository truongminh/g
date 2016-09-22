package s2

func (pq s2CellQueue) Len() int { return len(pq) }

func (pq s2CellQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	// the older the location, the higher the priority
	return pq[i].S2CellID() < pq[j].S2CellID()
}

func (pq s2CellQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].S2SetIndex(i)
	pq[j].S2SetIndex(j)
}

func (pq *s2CellQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(IEntry)
	item.S2SetIndex(n)
	*pq = append(*pq, item)
}

func (pq *s2CellQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	item.S2SetIndex(-1) // for safety
	return item
}
