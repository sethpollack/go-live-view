package liveview

import (
	"go-live-view/channel"
)

type Socket interface {
	channel.Socket
	PushEvent(string, any) error
	PushPatch(string, bool) error
	PushNavigate(string, bool) error
	Redirect(string) error
	Redirected() bool
	// PutFlash(string, string) error
}

var _ Socket = (*socket)(nil)

type socket struct {
	channel.Socket
	redirected bool
}

func newSocket(s channel.Socket) *socket {
	return &socket{
		Socket: s,
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
func (s *socket) PushPatch(url string, replace bool) error {
	err := s.Push("live_patch", map[string]any{
		"kind": navType(replace),
		"to":   url,
	})
	if err != nil {
		return err
	}

	// client does not return an event, so we push it to ourselves
	err = s.Socket.PushSelf("live_patch", map[string]any{
		"url": url,
	})
	if err != nil {
		return err
	}

	s.redirected = true

	return nil
}

// PushNavigate sends a live_redirect to the client.
func (s socket) PushNavigate(url string, replace bool) error {
	err := s.Push("live_redirect", map[string]any{
		"kind": navType(replace),
		"to":   url,
	})
	if err != nil {
		return err
	}

	s.redirected = true

	return nil
}

// Redirect sends a redirect to the client.
func (s socket) Redirect(url string) error {
	err := s.Push("redirect", map[string]any{
		"to": url,
	})
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

func navType(replace bool) string {
	if replace {
		return "replace"
	}
	return "push"
}
