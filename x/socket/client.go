package socket

import (
	"bytes"
	"encoding/json"
	"g/x/math"
)

type WsClient struct {
	ReadWriter
	UID  string
	Auth Auth
}

func (c *WsClient) Error(err error) {
	c.Write(BuildStringMessage("/error", err.Error()))
}

func (c *WsClient) WriteJson(uri string, v interface{}) {
	c.Write(BuildJsonMessage(uri, v))
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

var idMakerChanWsClient = math.RandStringMaker{Prefix: "chan", Length: 10}

const chanWriterLength = 10

type ChanWsClient struct {
	WsClient
	chanWriter
}

func (c *ChanWsClient) Close() {
	if c.chanWriter != nil {
		c.chanWriter.Close()
		c.chanWriter = nil
	}
}

type chanWriter chan []byte

func NewChanWsClient(a Auth) *WsClient {
	var c = &WsClient{
		Auth:       a,
		UID:        idMakerChanWsClient.Next(),
		ReadWriter: chanWriter(make(chan []byte, chanWriterLength)),
	}
	return c
}

func (c chanWriter) Write(data []byte) {
	c <- data
}

func (c chanWriter) Read() ([]byte, bool) {
	var bytes, ok = <-c
	return bytes, ok
}

func (c chanWriter) Close() {
	close(c)
}
