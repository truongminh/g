package socket

type EventHandler func(uri string, v interface{})
type IBoxHandler func(w ResponseWriter, r *Request)

type Auth interface {
	ID() string
}

type AuthID string

func (a *AuthID) ID() string {
	return string(*a)
}

var AuthOff = AuthID("")
