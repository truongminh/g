package socket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
)

func AddPrefixEventHandler(prefix string, handler EventHandler) EventHandler {
	return func(uri string, v interface{}) {
		handler(prefix+uri, v)
	}
}

type Request struct {
	Payload []byte
	RawURI  []byte
	URI     *url.URL
	Data    []byte
	Auth    Auth
}

func NewRequest(a Auth, payload []byte) *Request {
	var r = &Request{
		Auth:    a,
		Payload: payload,
	}

	var endOfURI = bytes.Index(payload, []byte(" "))
	var remaining = payload
	if endOfURI < 0 {
		r.RawURI = remaining
		remaining = remaining[0:0]
	} else {
		r.RawURI = remaining[:endOfURI]
		remaining = remaining[endOfURI+1:]
	}

	r.URI, _ = url.ParseRequestURI(string(r.RawURI))
	r.Data = remaining
	return r
}

func (r *Request) Path() string {
	return r.URI.Path
}

func (r *Request) UnmarshalJson(v interface{}) error {
	return json.Unmarshal(r.Data, v)
}

func (r *Request) String() string {
	return fmt.Sprintf("url: [%s], data: [%s]\n", r.RawURI, r.Data)
}
