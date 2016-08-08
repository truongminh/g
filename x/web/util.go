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

func MustGetInt64(key string, g IGetable) int64 {
	var v, err = strconv.ParseInt(g.Get(key), 10, 64)
	if err != nil {
		panic(BadRequest(key + " must be int"))
	}
	return v
}
