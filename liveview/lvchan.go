package liveview

import (
	"fmt"

	"github.com/sethpollack/go-live-view/channel"
	"github.com/sethpollack/go-live-view/params"
	"github.com/sethpollack/go-live-view/rend"
)

var _ channel.Channel = &lvChannel{}

type LvLifecycle interface {
	Join(Socket, params.Params) (*rend.Root, error)
	Leave() error
	StaticRender(string) (string, error)
	Event(Socket, params.Params) (*rend.Root, error)
	Params(Socket, params.Params) (*rend.Root, error)
	AllowUpload(Socket, params.Params) (any, error)
	Progress(Socket, params.Params) (*rend.Root, error)
	DestroyCIDs([]int) error
}

type lvChannel struct {
	lc LvLifecycle
}

func NewLVChannel(lc LvLifecycle) func() channel.Channel {
	return func() channel.Channel {
		return &lvChannel{
			lc: lc,
		}
	}
}

func (l *lvChannel) Join(s channel.Socket, p any) error {
	rend, err := l.lc.Join(newSocket(s), params.FromAny(p))
	if err != nil {
		return err
	}

	return s.Push("rendered", rend)
}

func (l *lvChannel) Leave(s channel.Socket) error {
	err := l.lc.Leave()
	if err != nil {
		return err
	}

	return s.Push("", nil)
}

func (l *lvChannel) Message(s channel.Socket, event string, p any) error {
	params := params.FromAny(p)

	switch event {
	case "live_patch":
		return l.handleLivePatchEvent(s, params)
	case "event":
		return l.handleEvent(s, params)
	case "allow_upload":
		return l.handleAllowUploadEvent(s, params)
	case "progress":
		return l.handleProgressEvent(s, params)
	case "cids_destroyed", "cids_will_destroy":
		return l.handleDestroyCidsEvent(s, params)
	default:
		return fmt.Errorf("unhandled event: %s", event)
	}
}

func (l *lvChannel) Broadcast(s channel.Socket, event string, p any) error {
	params := params.FromAny(p)

	switch event {
	case "event":
		return l.handleEvent(s, params)
	case "live_patch":
		return l.handleLivePatchEvent(s, params)
	default:
		return fmt.Errorf("unhandled event: %s", event)
	}
}

func (l *lvChannel) handleEvent(s channel.Socket, p params.Params) error {
	diff, err := l.lc.Event(newSocket(s), p)
	if err != nil {
		return err
	}
	return s.Push("diff", diff)
}

func (l *lvChannel) handleDestroyCidsEvent(s channel.Socket, p params.Params) error {
	cids := p.IntSlice("cids")

	err := l.lc.DestroyCIDs(cids)
	if err != nil {
		return err
	}

	return s.Push("", cids)
}

func (l *lvChannel) handleLivePatchEvent(s channel.Socket, p params.Params) error {
	diff, err := l.lc.Params(newSocket(s), p)
	if err != nil {
		return err
	}

	return s.Push("diff", diff)
}

func (l *lvChannel) handleAllowUploadEvent(s channel.Socket, p params.Params) error {
	payload, err := l.lc.AllowUpload(newSocket(s), p)
	if err != nil {
		return err
	}

	return s.Push("", payload)
}

func (l *lvChannel) handleProgressEvent(s channel.Socket, p params.Params) error {
	payload, err := l.lc.Progress(newSocket(s), p)
	if err != nil {
		return err
	}

	return s.Push("diff", payload)
}
