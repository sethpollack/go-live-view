package channel

import (
	"context"
	"fmt"
	"strings"
	"sync"
)

type server struct {
	mu sync.RWMutex

	h        *Hub
	c        *conn
	matchers map[string]func() Channel
	channels map[string]Channel
}

func NewServer(c Conn, h *Hub) *server {
	return &server{
		h:        h,
		c:        newConnection(c),
		matchers: make(map[string]func() Channel),
		channels: make(map[string]Channel),
	}
}

func (s *server) Route(topic string, factory func() Channel) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.matchers[topic] = factory
}

func (s *server) Close(topic string) {
	s.deleteChannel(topic)
}

func (s *server) Broadcast(msg *Message) error {
	mChan, err := s.getChannel(msg.Topic)
	if err != nil {
		return err
	}

	sock := NewSocket(s, msg)

	return mChan.Broadcast(sock, msg.Event, msg.Payload)
}

func (s *server) PushBroadcast(msg *Message) error {
	if s.h == nil {
		return fmt.Errorf("no server available")
	}

	return s.h.WriteMessage(msg)
}

func (s *server) Push(msg *Message) error {
	return s.c.WriteMessage(msg)
}

func (s *server) Listen(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			msg, err := s.c.ReadMessage()
			if err != nil {
				return
			}

			switch msg.Event {
			case "heartbeat":
				err := s.handleHeartbeat(msg)
				if err != nil {
					s.handleError(msg, err)
				}
			case "phx_join":
				err := s.handleJoin(msg)
				if err != nil {
					s.handleError(msg, err)
				}
			case "phx_leave":
				err := s.handleLeave(msg)
				if err != nil {
					s.handleError(msg, err)
				}
			default:
				err := s.handleMessage(msg)
				if err != nil {
					s.handleError(msg, err)
				}
			}

		}
	}
}

func (s *server) handleHeartbeat(msg *Message) error {
	return s.Push(&Message{
		JoinRef: msg.JoinRef,
		Ref:     msg.Ref,
		Topic:   msg.Topic,
		Event:   "phx_reply",
		Payload: map[string]any{
			"status":   "ok",
			"response": map[string]any{},
		},
	})
}

func (s *server) handleJoin(msg *Message) error {
	mChan, err := s.match(msg.Topic)
	if err != nil {
		return err
	}

	sock := NewSocket(s, msg)

	err = mChan.Join(sock, msg.Payload)
	if err != nil {
		return err
	}

	s.setChannel(msg.Topic, mChan)

	return nil
}

func (s *server) handleLeave(msg *Message) error {
	mChan, err := s.getChannel(msg.Topic)
	if err != nil {
		return err
	}

	sock := NewSocket(s, msg)

	err = mChan.Leave(sock)
	if err != nil {
		return err
	}

	s.deleteChannel(msg.Topic)

	return nil
}

func (s *server) handleMessage(msg *Message) error {
	mChan, err := s.getChannel(msg.Topic)
	if err != nil {
		return err
	}

	sock := NewSocket(s, msg)

	return mChan.Message(sock, msg.Event, msg.Payload)
}

func (s *server) handleError(msg *Message, err error) {
	pushErr := s.Push(&Message{
		JoinRef: msg.JoinRef,
		Ref:     msg.Ref,
		Topic:   msg.Topic,
		Event:   "phx_reply",
		Payload: map[string]any{
			"status": "error",
			"response": map[string]any{
				"reason": err.Error(),
			},
		},
	})
	if pushErr != nil {
		fmt.Println(pushErr)
	}
}

func (s *server) match(topic string) (Channel, error) {
	for t, f := range s.matchers {
		if match(t, topic) {
			return f(), nil
		}
	}

	return nil, fmt.Errorf("no channel found for topic %s", topic)
}

func (s *server) getChannel(topic string) (Channel, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	channel, ok := s.channels[topic]
	if !ok {
		return nil, fmt.Errorf("no channel found for topic %s", topic)
	}

	return channel, nil
}

func (s *server) setChannel(topic string, c Channel) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.channels[topic] = c
}

func (s *server) deleteChannel(topic string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.channels, topic)
}

func match(pattern, path string) bool {
	p1 := strings.Split(pattern, ":")
	p2 := strings.Split(path, ":")

	if len(p1) != len(p2) {
		return false
	}

	if p1[0] != p2[0] {
		return false
	}

	if p1[1] == "*" {
		return true
	}

	if p1[1] != p2[1] {
		return false
	}

	return true
}
