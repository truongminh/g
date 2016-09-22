package s2

import (
	"github.com/golang/geo/s2"
)

// The Earth's mean radius in kilometers (according to NASA).
const earthRadiusKm float32 = 6371.01

var km20Coverer = &s2.RegionCoverer{
	MinLevel: 20,
	MaxCells: 32,
}

func goodCoverer(km float32) *s2.RegionCoverer {
	return km20Coverer
}

func distance(c1 s2.CellID, c2 s2.CellID) float32 {
	return earthRadiusKm * float32(c1.Point().Distance(c2.Point()))
}
