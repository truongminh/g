package socket

type SubscriptionManagement struct {
	subscribers map[string][]ResponseWriter
}

func newSubscriptionManagement() *SubscriptionManagement {
	return &SubscriptionManagement{
		subscribers: make(map[string][]ResponseWriter),
	}
}

func (s *SubscriptionManagement) Line(uri string) []ResponseWriter {
	if s.subscribers[uri] == nil {
		s.subscribers[uri] = make([]ResponseWriter, 0)
	}
	return s.subscribers[uri]
}

func (s *SubscriptionManagement) Subscribe(w ResponseWriter, uri string) {
	var subscribers = s.Line(uri)
	s.subscribers[uri] = append(subscribers, w)
}

func (s *SubscriptionManagement) Unsubscribe(w ResponseWriter) {

}