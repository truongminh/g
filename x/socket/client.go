package socket

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type WsClient interface {
	Write(p []byte)
	UID() string
	Auth() Auth
}

func SendError(w WsClient, err error) {
	w.Write(BuildStringMessage("/error", err.Error()))
}

func SendJson(w WsClient, uri string, v interface{}) {
	w.Write(BuildJsonMessage(uri, v))
}

func Send(w WsClient, payload []byte) {
	w.Write(payload)
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

type ChanWsClient struct {
	send chan []byte
	id   string
	auth Auth
}

func NewChanWsClient(a Auth) *ChanWsClient {
	var c = &ChanWsClient{}
	c.send = make(chan []byte)
	c.id = fmt.Sprintf("chan-%v", c)
	c.auth = a
	return c
}

func (c *ChanWsClient) Write(data []byte) {
	c.send <- data
}

func (c *ChanWsClient) UID() string {
	return c.id
}

func (c *ChanWsClient) Auth() Auth {
	return c.auth
}
