package socket

import (
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

	var c = newWsClient(a, ws)
	b.Clients.Add(c, c.Auth.ID())
	if b.Clients != nil {
		for _, item := range b.Clients {
			var te = item
			glog.Info(te)
		}
	}

	if onJoin != nil {
		onJoin(c)
	}
	var done = make(chan struct{})

	defer func() {
		b.Clients.Remove(c, c.Auth.ID())
		c.Close()
		<-done
		if onLeave != nil {
			onLeave(c)
		}
	}()

	go func() {
		for {
			var bytes, ok = <-c.reply
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
			glog.Info(err)
			break
		}
		var r, err = NewRequest(c, data)
		if err != nil {
			c.WriteError(err)
		} else {
			b.Serve(r)
		}
	}
}

var (
	errHandlerNotFound = BadRequest("HANDLER NOT FOUND")
	errInternalServer  = InternalServerError("SERVER ERROR")
)

func (b *Box) notFound(r *Request) {
	r.ReplyError(errHandlerNotFound)
}

func (b *Box) defaultRecover(r *Request, rc interface{}) {
	if err, ok := rc.(error); ok {
		if _, ok = err.(IWebError); ok {
			r.ReplyError(err)
			return
		}
		glog.Error(err, string(debug.Stack()))
		r.ReplyError(errInternalServer)
	} else {
		glog.Error(rc, string(debug.Stack()))
		r.ReplyError(errInternalServer)
	}
}

func (b *Box) join(w *WsClient) {

}

func (b *Box) leave(w *WsClient) {

}

type WsServer struct{}

func (s *WsServer) WriteError(ws *websocket.Conn, err error) {
	websocket.Message.Send(ws, string(BuildErrorMessage("/system", err)))
}

func (s *WsServer) Recover(ws *websocket.Conn) {
	if r := recover(); r != nil {
		if err, ok := r.(error); ok {
			if _, ok = err.(IWebError); ok {
				s.WriteError(ws, err)
				return
			}
			glog.Error(err, string(debug.Stack()))
			s.WriteError(ws, errInternalServer)
		} else {
			glog.Error(r, string(debug.Stack()))
			s.WriteError(ws, errInternalServer)
		}
	}
}

func (b *Box) WriteError(ws *websocket.Conn, err error) {
	websocket.Message.Send(ws, string(BuildErrorMessage("/system", err)))
}
