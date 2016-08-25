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
	RawURI  string
	URI     *url.URL
	Data    []byte
	Auth    Auth
}

func NewRequest(a Auth, payload []byte) (*Request, error) {
	var r = &Request{
		Auth:    a,
		Payload: payload,
	}

	var endOfURI = bytes.Index(payload, []byte(" "))
	var remaining = payload
	if endOfURI < 0 {
		r.RawURI = string(remaining)
		remaining = remaining[0:0]
	} else {
		r.RawURI = string(remaining[:endOfURI])
		remaining = remaining[endOfURI+1:]
	}

	var err error

	r.URI, err = url.ParseRequestURI(r.RawURI)
	if err != nil {
		return nil, err
	}
	r.Data = remaining
	return r, nil
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

func (r *Request) Reply(w ResponseWriter, v interface{}) {
	SendJson(w, r.RawURI, v)
}

func (r *Request) ReplyError(w ResponseWriter, v interface{}) {
	SendJson(w, "/error/"+r.RawURI, v)
}
