package channel

import (
	"context"
	"sync"
)

type broadcaster interface {
	Broadcast(msg *Message) error
}

type hub struct {
	mu      sync.RWMutex
	servers map[broadcaster]struct{}
}

func NewHub() *hub {
	return &hub{
		servers: make(map[broadcaster]struct{}),
	}
}

func (sr *hub) Add(s broadcaster) {
	sr.mu.Lock()
	defer sr.mu.Unlock()

	sr.servers[s] = struct{}{}
}

func (sr *hub) Remove(s broadcaster) {
	sr.mu.Lock()
	defer sr.mu.Unlock()

	delete(sr.servers, s)
}

func (sr *hub) WriteMessage(msg *Message) error {
	// send to all sockets
	for s := range sr.servers {
		err := s.Broadcast(msg)
		if err != nil {
			return err
		}
	}
	// sr.Broadcast(msg) // TODO: send to pubsub
	return nil
}

func (sr *hub) Listen(ctx context.Context) {
	// TODO add listener for pubsub
	for {
		select {
		case <-ctx.Done():
			return
		}
	}
}
