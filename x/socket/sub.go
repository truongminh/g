package socket

type SubscriptionManagement struct {
	subscribers map[string][]WsClient
}

func newSubscriptionManagement() *SubscriptionManagement {
	return &SubscriptionManagement{
		subscribers: make(map[string][]WsClient),
	}
}

func (s *SubscriptionManagement) Line(uri string) []WsClient {
	if s.subscribers[uri] == nil {
		s.subscribers[uri] = make([]WsClient, 0)
	}
	return s.subscribers[uri]
}

func (s *SubscriptionManagement) Subscribe(w WsClient, uri string) {
	var subscribers = s.Line(uri)
	s.subscribers[uri] = append(subscribers, w)
}

func (s *SubscriptionManagement) Unsubscribe(w WsClient) {
	for uri, line := range s.subscribers {
		var newLine = make([]WsClient, 0)
		for _, sw := range line {
			if sw.UID() != w.UID() {
				newLine = append(newLine, sw)
			} else {
				//
			}
		}
		s.subscribers[uri] = newLine
	}
}
