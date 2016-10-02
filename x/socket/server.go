package socket

import (
	"errors"
	"runtime/debug"

	"github.com/golang/glog"
	"golang.org/x/net/websocket"
)

func (b *Box) AcceptPublic(ws *websocket.Conn, args ...Auth) {
	if len(args) < 1 {
		b.AcceptDefault(ws, AuthOff)
	} else {
		b.AcceptDefault(ws, args[0])
	}
}

func (b *Box) AcceptDefault(ws *websocket.Conn, a Auth) {
	b.Accept(ws, a, b.Join, b.Leave)
}

func (b *Box) Accept(ws *websocket.Conn, a Auth, onJoin func(*WsClient), onLeave func(*WsClient)) {

	var codec = websocket.Message

	var c = NewChanWsClient(a)
	b.Clients.Add(c, c.Auth.ID())
	if onJoin != nil {
		onJoin(c)
	}
	var done = make(chan struct{})

	defer func() {
		c.Close()
		<-done
		b.Clients.Remove(c, c.Auth.ID())
		if onLeave != nil {
			onLeave(c)
		}
	}()

	go func() {
		for {
			var bytes, ok = c.Read()
			if !ok {
				break
			}
			var data = string(bytes)
			if err := codec.Send(ws, data); err != nil {
				break
			}
		}
		done <- struct{}{}
	}()

	for {
		var data []byte
		if err := codec.Receive(ws, &data); err != nil {
			break
		}
		var r, err = NewRequest(c, data)
		if err != nil {
			c.Error(err)
		} else {
			b.Serve(r)
		}
	}
}

var (
	errHandlerNotFound = errors.New("HANDLER NOT FOUND")
	errInternalServer  = errors.New("SERVER ERROR")
)

func (b *Box) notFound(r *Request) {
	r.Error(errHandlerNotFound)
}

func (b *Box) defaultRecover(r *Request, rc interface{}) {
	if err, ok := rc.(error); ok {
		if _, ok = err.(IWebError); ok {
			r.Error(err)
			return
		}
		glog.Error(err, string(debug.Stack()))
		r.Error(errInternalServer)
	} else {
		glog.Error(rc, string(debug.Stack()))
		r.Error(errInternalServer)
	}
}

func (b *Box) join(w *WsClient) {

}

func (b *Box) leave(w *WsClient) {

}
