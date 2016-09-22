package s2

import (
	"github.com/golang/geo/s2"
)

func (m *Map) ForEachWithin(lat float32, lng float32, km float32, f func(IEntry, float32)) {
	var dOnEarth = float64(km / earthRadiusKm)
	var searchRegion = s2.RectFromCenterSize(
		s2.LatLngFromDegrees(float64(lat), float64(lng)),
		s2.LatLngFromDegrees(dOnEarth, dOnEarth),
	)

	var coverer = goodCoverer(km)
	var cells = coverer.Covering(searchRegion)
	var searchCell = MakeCell(lat, lng)

	for _, cell := range cells {
		m.q.ascendRange(cell.RangeMin(), cell.RangeMax(), func(loc IEntry) {
			var d = distance(searchCell, loc.S2CellID())
			if d < km {
				f(loc, d)
			}
		})
	}
}
