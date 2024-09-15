package liveview

import (
	"encoding/base64"
	"encoding/json"

	"github.com/sethpollack/go-live-view/channel"
)

var _ Socket = (*socket)(nil)

type redirectOption func(map[string]any)

type Socket interface {
	channel.Socket
	PushEvent(string, any) error
	PushPatch(string, ...redirectOption) error
	PushNavigate(string, ...redirectOption) error
	Redirect(string, ...redirectOption) error
	Redirected() bool
}

type socket struct {
	channel.Socket
	redirected bool
}

func NewSocket(s channel.Socket) *socket {
	return &socket{
		Socket: s,
	}
}

func WithFlash(key, value string) redirectOption {
	return func(m map[string]any) {
		m["flash"] = map[string]string{
			key: value,
		}
	}
}

func WithReplace() redirectOption {
	return func(m map[string]any) {
		m["kind"] = "replace"
	}
}

// PushSelf sends payload back to the mounted liveview.
func (s *socket) PushSelf(event string, payload any) error {
	return s.Socket.PushSelf("event",
		map[string]any{
			"event": event,
			"type":  "self",
			"value": payload,
		},
	)
}

// PushBroadcast sends payload to all liveviews.
func (s *socket) PushBroadcast(event string, payload any) error {
	return s.Socket.PushBroadcast("event",
		map[string]any{
			"event": event,
			"type":  "broadcast",
			"value": payload,
		},
	)
}

// PushEvent sends an event to the client.
func (s *socket) PushEvent(event string, payload any) error {
	return s.Push("e", [][]any{
		{
			event, payload,
		},
	})
}

// PushPatch sends a live_patch to the client.
func (s *socket) PushPatch(url string, opts ...redirectOption) error {
	payload := map[string]any{
		"to":   url,
		"kind": "push",
	}

	for _, opt := range opts {
		opt(payload)
	}

	err := s.Push("live_patch", payload)
	if err != nil {
		return err
	}

	// client does not return an event, so we push it to ourselves
	err = s.Socket.PushSelf("live_patch", map[string]any{
		"kind":  payload["kind"],
		"url":   payload["to"],
		"flash": payload["flash"],
	})
	if err != nil {
		return err
	}

	s.redirected = true

	return nil
}

// PushNavigate sends a live_redirect to the client.
func (s socket) PushNavigate(url string, opts ...redirectOption) error {
	payload := map[string]any{
		"to":   url,
		"kind": "push",
	}

	for _, opt := range opts {
		opt(payload)
	}

	err := s.Push("live_redirect", payload)
	if err != nil {
		return err
	}

	s.redirected = true

	return nil
}

// Redirect sends a redirect to the client.
func (s socket) Redirect(url string, opts ...redirectOption) error {
	payload := map[string]any{
		"to": url,
	}

	for _, opt := range opts {
		opt(payload)
	}

	encodeFlash(payload)

	err := s.Push("redirect", payload)
	if err != nil {
		return err
	}

	s.redirected = true

	return nil
}

// Redirected returns true if the socket has been redirected.
func (s *socket) Redirected() bool {
	return s.redirected
}

func encodeFlash(m map[string]any) {
	flash, ok := m["flash"]
	if !ok {
		return
	}
	delete(m, "flash")

	js, err := json.Marshal(flash)
	if err != nil {
		return
	}

	m["flash"] = base64.StdEncoding.EncodeToString(js)
}
