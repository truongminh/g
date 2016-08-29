package socket

import (
	"errors"
	"github.com/golang/glog"
	"golang.org/x/net/websocket"
	"runtime/debug"
	"sync"
)

func (b *Box) AcceptPublic(ws *websocket.Conn, args ...Auth) {
	if len(args) < 1 {
		b.Accept(ws, &AuthOff)
	} else {
		b.Accept(ws, args[0])
	}
}

func (b *Box) Accept(ws *websocket.Conn, a Auth) {

	var wait = sync.WaitGroup{}
	var codec = websocket.Message

	var w = NewChanResponseWriter()
	b.Join(w, a)

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
	close(w.send)
	wait.Wait()
	b.SubManager.Unsubscribe(w)
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
