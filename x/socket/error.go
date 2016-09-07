package socket

type IWebError interface {
	StatusCode() int
}

func AssertNil(err error) {
	if err != nil {
		panic(err)
	}
}
