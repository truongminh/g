package s2

import (
	"sort"

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

type view struct {
	IEntry
	Km float32
}

type viewByKm []view

func (a viewByKm) Len() int           { return len(a) }
func (a viewByKm) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a viewByKm) Less(i, j int) bool { return a[i].Km < a[j].Km }

func (m *Map) NearestWithin(lat float32, lng float32, km float32, limit int32) []view {
	var res = viewByKm(make([]view, 0))
	m.ForEachWithin(lat, lng, km, func(e IEntry, d float32) {
		res = append(res, view{IEntry: e, Km: d})
	})
	sort.Sort(res)
	return res
}
