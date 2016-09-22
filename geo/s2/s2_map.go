package s2

type Map struct {
	q *s2CellQueue
}

func NewMap() *Map {
	var q = s2CellQueue([]IEntry{})
	return &Map{q: &q}
}
