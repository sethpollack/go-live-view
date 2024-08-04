package channel

type Socket interface {
	Push(string, any) error
	PushBroadcast(string, any) error
	PushSelf(string, any) error
	Close() error
}

type socket struct {
	server  *server
	joinRef string
	ref     string
	topic   string
	replied bool
}

func NewSocket(s *server, msg *Message) *socket {
	return &socket{
		server:  s,
		joinRef: msg.JoinRef,
		ref:     msg.Ref,
		topic:   msg.Topic,
	}
}

func (s *socket) Push(event string, payload any) error {
	if payload == nil {
		payload = map[string]any{}
	}

	if s.ref != "" && event != "" {
		payload = map[string]any{
			"status": "ok",
			"response": map[string]any{
				event: payload,
			},
		}
		event = "phx_reply"
	}

	if s.ref != "" && event == "" {
		payload = map[string]any{
			"status":   "ok",
			"response": payload,
		}
		event = "phx_reply"
	}

	// prevent multiple replies
	if event == "phx_reply" {
		if s.replied {
			return nil
		}
		s.replied = true
	}

	return s.server.Push(&Message{
		JoinRef: s.joinRef,
		Ref:     s.ref,
		Topic:   s.topic,
		Event:   event,
		Payload: payload,
	})
}

func (s *socket) Close() error {
	s.server.Close(s.topic)

	return s.server.Push(&Message{
		JoinRef: s.joinRef,
		Topic:   s.topic,
		Event:   "phx_close",
		Payload: map[string]any{},
	})
}

func (s *socket) PushSelf(event string, payload any) error {
	return s.server.Broadcast(&Message{
		JoinRef: s.joinRef,
		Topic:   s.topic,
		Event:   event,
		Payload: payload,
	})
}

func (s *socket) PushBroadcast(event string, payload any) error {
	return s.server.PushBroadcast(&Message{
		JoinRef: s.joinRef,
		Topic:   s.topic,
		Event:   event,
		Payload: payload,
	})
}
