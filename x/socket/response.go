package socket

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type ResponseWriter interface {
	Write(p []byte)
	SubscribeAs() string
}

func SendError(w ResponseWriter, err error) {
	w.Write(BuildStringMessage("/error", err.Error()))
}

func SendJson(w ResponseWriter, uri string, v interface{}) {
	w.Write(BuildJsonMessage(uri, v))
}

func BuildJsonMessage(uri string, v interface{}) []byte {
	var data, _ = json.Marshal(v)
	return BuildRawMessage([]byte(uri), data)
}

func BuildStringMessage(uri string, v string) []byte {
	return BuildRawMessage([]byte(uri), []byte(v))
}

func BuildRawMessage(uri []byte, data []byte) []byte {
	var buffer = bytes.NewBuffer(uri)
	buffer.WriteString(" ")
	buffer.Write(data)
	return buffer.Bytes()
}

type ChanResponseWriter struct {
	send chan []byte
	id   string
}

func NewChanResponseWriter() *ChanResponseWriter {
	var c = &ChanResponseWriter{}
	c.send = make(chan []byte)
	c.id = fmt.Sprintf("chan-%v", c)
	return c
}

func (c *ChanResponseWriter) Write(data []byte) {
	c.send <- data
}

func (c *ChanResponseWriter) SubscribeAs() string {
	return c.id
}
