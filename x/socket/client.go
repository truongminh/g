package socket

import (
	"g/x/math"

	"golang.org/x/net/websocket"
)

type WsClient struct {
	UID    string
	Auth   Auth
	Socket *websocket.Conn
	reply  chan []byte
}

func (c *WsClient) write(data []byte) {
	c.reply <- data
}

func (c *WsClient) WriteError(err error) {
	c.reply <- BuildErrorMessage("/server", err)
}

func (c *WsClient) WriteJson(uri string, v interface{}) {
	c.reply <- BuildJsonMessage(uri, v)
}

var idMakerChanWsClient = math.RandStringMaker{Prefix: "chan", Length: 10}

const chanWriterLength = 64

func newWsClient(a Auth, s *websocket.Conn) *WsClient {
	var c = &WsClient{
		Auth:   a,
		Socket: s,
		UID:    idMakerChanWsClient.Next(),
		reply:  make(chan []byte, chanWriterLength),
	}
	return c
}

func (c *WsClient) Close() {
	if c.Socket != nil {
		c.Socket.Close()
		close(c.reply)
		c.Socket = nil
	}
}
