package web

import (
	"encoding/json"
	"github.com/golang/glog"
	"net/http"
	"runtime/debug"
)

type JsonServer struct{}

func (s *JsonServer) SendError(w http.ResponseWriter, err error) {
	if err != nil {
		if werr, ok := err.(IWebError); ok {
			w.WriteHeader(werr.StatusCode())
		} else {
			glog.Error(debug.Stack())
			err = ErrServerError
			w.WriteHeader(http.StatusInternalServerError)
		}
		s.SendJson(w, map[string]string{
			"status": "error",
			"error":  err.Error(),
		})
	}
}

func (s *JsonServer) SendJson(w http.ResponseWriter, v interface{}) {
	json.NewEncoder(w).Encode(v)
}

func (s *JsonServer) SendData(w http.ResponseWriter, v interface{}) {
	s.SendJson(w, map[string]interface{}{
		"status": "success",
		"data":   v,
	})
}

func (s *JsonServer) Success(w http.ResponseWriter) {
	s.SendData(w, nil)
}

func (s *JsonServer) DecodeBody(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func (s *JsonServer) MustDecodeBody(r *http.Request, v interface{}) {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		panic(BadRequest(err.Error()))
	}
}

func (s *JsonServer) Recover(w http.ResponseWriter) {
	if r := recover(); r != nil {
		if err, ok := r.(error); ok {
			s.SendError(w, err)
		} else {
			s.SendError(w, ErrServerError)
			glog.Error(debug.Stack())
		}
	}
}