package web

import (
	"strconv"
)

type IGetable interface {
	Get(key string) string
}

func ParseFloat64(key string, g IGetable) (float64, error) {
	var v, err = strconv.ParseFloat(g.Get(key), 64)
	if err != nil {
		return 0, BadRequest(key + " must be float64")
	}
	return v, nil
}
