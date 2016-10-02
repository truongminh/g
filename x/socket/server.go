package socket

import (
	"errors"
	"runtime/debug"
	"sync"

	"github.com/golang/glog"
	"golang.org/x/net/websocket"
)

func (b *Box) AcceptPublic(ws *websocket.Conn, args ...Auth) {
	if len(args) < 1 {
		b.Accept(ws, AuthOff)
	} else {
		b.Accept(ws, args[0])
	}
}

func (b *Box) Accept(ws *websocket.Conn, a Auth) {

	var wait = sync.WaitGroup{}
	var codec = websocket.Message

	var w = NewChanWsClient(a)

	defer func() {
		close(w.send)
		wait.Wait()
		b.Clients.Remove(w)
		b.SubManager.Unsubscribe(w)
	}()

	go func() {
		wait.Add(1)
		for {
			var bytes, ok = <-w.send
			if !ok {
				break
			}
			var data = string(bytes)
			if err := codec.Send(ws, data); err != nil {
				break
			}
		}
		wait.Done()
	}()

	b.Join(w)
	b.Clients.Add(w)

	for {
		var data []byte
		if err := codec.Receive(ws, &data); err != nil {
			break
		}
		var r, err = NewRequest(w, data)
		if err != nil {
			SendError(w, err)
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

func (b *Box) join(w WsClient) {

}
