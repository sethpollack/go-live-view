package liveview

import (
	"github.com/sethpollack/go-live-view/params"
	"github.com/sethpollack/go-live-view/rend"
	"github.com/sethpollack/go-live-view/uploads"
)

type View interface {
	Render(rend.Node) (rend.Node, error)
}

type Mounter interface {
	Mount(Socket, params.Params) error
}

type Unmounter interface {
	Unmount() error
}

type Patcher interface {
	Params(Socket, params.Params) error
}

type EventHandler interface {
	Event(Socket, string, params.Params) error
}

type Uploader interface {
	Uploads() *uploads.Uploads
}

func TryMount(a any, s Socket, p params.Params) error {
	if m, ok := a.(Mounter); ok {
		return m.Mount(s, p)
	}

	return nil
}

func TryUnmount(a any) error {
	if m, ok := a.(Unmounter); ok {
		return m.Unmount()
	}

	return nil
}

func TryParams(a any, s Socket, p params.Params) error {
	if m, ok := a.(Patcher); ok {
		return m.Params(s, p)
	}

	return nil
}

func TryEvent(a any, s Socket, event string, p params.Params) error {
	if m, ok := a.(EventHandler); ok {
		return m.Event(s, event, p)
	}

	return nil
}

func TryUploads(a any) *uploads.Uploads {
	if m, ok := a.(Uploader); ok {
		return m.Uploads()
	}

	return nil
}
