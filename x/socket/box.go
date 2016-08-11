package socket

// Box handles ws request
//
type Box struct {
	ID          string
	Subscribers map[string]ResponseWriter
	SubManager  *SubscriptionManagement
	handlers    map[string]IBoxHandler
	NotFound    IBoxHandler
	Join        func(ResponseWriter, Auth)
	Recover     func(ResponseWriter, *Request, interface{})
}

// NewBox create a new box
func NewBox(ID string) *Box {
	var b = &Box{
		ID:          ID,
		Subscribers: make(map[string]ResponseWriter),
		SubManager:  newSubscriptionManagement(),
		handlers:    make(map[string]IBoxHandler),
	}
	b.Recover = b.defaultRecover
	b.NotFound = b.notFound
	b.Join = b.join
	b.Handle("/echo", b.Echo)
	return b
}

// Handle add a handler
func (b *Box) Handle(uri string, handler IBoxHandler) {
	b.handlers[uri] = handler
}

// Serve process the request
func (b *Box) Serve(w ResponseWriter, r *Request) {

	defer func() {
		if rc := recover(); rc != nil {
			b.Recover(w, r, rc)
		}
	}()

	var handler = b.handlers[r.Path()]
	if handler == nil {
		handler = b.NotFound
	}
	handler(w, r)
}

// Echo the default echo service
func (b *Box) Echo(w ResponseWriter, r *Request) {
	w.Write(r.Payload)
}

// Emit send a message to some subscriber
func (b *Box) Emit(uri string, v interface{}) {
	var buffer = BuildJsonMessage(uri, v)
	for _, w := range b.SubManager.Line(uri) {
		w.Write(buffer)
	}
}
