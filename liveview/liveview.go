package liveview

import (
	"github.com/sethpollack/go-live-view/params"
	"github.com/sethpollack/go-live-view/rend"
	"github.com/sethpollack/go-live-view/uploads"
)

type LiveView interface {
	Mount(Socket, params.Params) error
	Unmount() error

	Params(Socket, params.Params) error
	Event(Socket, string, params.Params) error
	Render(rend.Node) (rend.Node, error)

	Uploads() *uploads.Uploads
}

type Base struct {
}

func NewBase() *Base {
	return &Base{}
}

var _ LiveView = &Base{}

func (b *Base) Mount(Socket, params.Params) error {
	return nil
}

func (b *Base) Unmount() error {
	return nil
}

func (b *Base) Params(Socket, params.Params) error {
	return nil
}

func (b *Base) Event(Socket, string, params.Params) error {
	return nil
}

func (b *Base) Render(rend.Node) (rend.Node, error) {
	return nil, nil
}

func (b *Base) Uploads() *uploads.Uploads {
	return nil
}
