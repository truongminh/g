package ex

// A expireQueue implements heap.Interface and holds Items.
type expireQueue []IEntry

func (exQ expireQueue) Len() int { return len(exQ) }

func (exQ expireQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	// the older the location, the higher the priority
	return exQ[i].ExMTime() > exQ[j].ExMTime()
}

func (exQ expireQueue) Swap(i, j int) {
	exQ[i], exQ[j] = exQ[j], exQ[i]
	exQ[i].ExSetIndex(i)
	exQ[j].ExSetIndex(j)
}

func (exQ *expireQueue) Push(x interface{}) {
	n := len(*exQ)
	item := x.(IEntry)
	item.ExSetIndex(n)
	*exQ = append(*exQ, item)
}

func (exQ *expireQueue) Pop() interface{} {
	old := *exQ
	n := len(old)
	item := old[n-1]
	item.ExSetIndex(-1)
	*exQ = old[0 : n-1]
	return item
}
