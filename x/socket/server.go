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

	var w = NewChanResponseWriter()

	defer func() {
		close(w.send)
		wait.Wait()
		delete(b.Writers, w.SubscribeAs())
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

	b.Join(w, a)
	b.Writers[w.SubscribeAs()] = w

	for {
		var data []byte
		if err := codec.Receive(ws, &data); err != nil {
			break
		}
		var r, err = NewRequest(a, data)
		if err != nil {
			SendError(w, err)
		} else {
			b.Serve(w, r)
		}
	}
}

func (b *Box) notFound(w ResponseWriter, request *Request) {
	SendError(w, errors.New("HANDLER NOT FOUND"))
}

func (b *Box) defaultRecover(w ResponseWriter, r *Request, rc interface{}) {
	if err, ok := rc.(error); ok {
		glog.Error(err, string(debug.Stack()))
		err = errors.New("server error")
		SendError(w, err)
	} else {
		glog.Error(rc, string(debug.Stack()))
		SendError(w, errors.New("server error"))
	}
}

func (b *Box) join(w ResponseWriter, a Auth) {

}

func (b *Box) Broadcast(uri string, v interface{}) {
	for _, w := range b.Writers {
		SendJson(w, uri, v)
	}
}
