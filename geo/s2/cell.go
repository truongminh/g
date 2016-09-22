package s2

import (
	"fmt"
	"github.com/golang/geo/s2"
)

type IEntry interface {
	S2GetIndex() int
	S2SetIndex(i int)
	S2CellID() s2.CellID
}

type S2Entry struct {
	Lat     float32
	Lng     float32
	S2Index int
	Cell    s2.CellID
}

func (e *S2Entry) S2GetIndex() int {
	return e.S2Index
}

func (e *S2Entry) S2SetIndex(i int) {
	e.S2Index = i
}

func (e *S2Entry) S2CellID() s2.CellID {
	return e.Cell
}

func (e *S2Entry) S2Move(lat float32, lng float32) {
	e.Lat = lat
	e.Lng = lng
	e.Cell = MakeCell(lat, lng)
}

// A s2CellQueue implements heap.Interface and holds Items.
type s2CellQueue []IEntry

func (pq s2CellQueue) Print(tag string) {
	fmt.Printf("\n%v\n", tag)
	for _, e := range pq {
		fmt.Printf("(%v) on heap %v \n", e.S2CellID().LatLng(), e.S2GetIndex())
	}
}

func MakeCell(lat float32, lng float32) s2.CellID {
	return s2.CellIDFromLatLng(s2.LatLngFromDegrees(float64(lat), float64(lng)))
}
