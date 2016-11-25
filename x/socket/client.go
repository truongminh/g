package socket

import "g/x/math"

type WsClient struct {
	ReadWriter
	UID  string
	Auth Auth
}

func (c *WsClient) WriteError(err error) {
	c.Write(BuildErrorMessage("/server", err))
}

func (c *WsClient) WriteJson(uri string, v interface{}) {
	c.Write(BuildJsonMessage(uri, v))
}

var idMakerChanWsClient = math.RandStringMaker{Prefix: "chan", Length: 10}

const chanWriterLength = 64

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
