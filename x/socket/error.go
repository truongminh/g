package socket

import (
	"net/http"
)

type IWebError interface {
	StatusCode() int
}

func AssertNil(err error) {
	if err != nil {
		panic(err)
	}
}

type BadRequest string

func (b BadRequest) StatusCode() int {
	return http.StatusBadRequest
}

func (b BadRequest) Error() string {
	return string(b)
}

func WrapBadRequest(err error, data string) error {
	if err != nil {
		return BadRequest(data + " " + err.Error())
	}
	return nil
}
