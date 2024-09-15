package lvchan

import (
	lv "github.com/sethpollack/go-live-view/liveview"
	"github.com/sethpollack/go-live-view/params"
	"github.com/sethpollack/go-live-view/rend"
)

type lifecycle interface {
	Join(lv.Socket, params.Params) (*rend.Root, error)
	Leave() error
	StaticRender(string) (string, error)
	Event(lv.Socket, params.Params) (*rend.Root, error)
	Params(lv.Socket, params.Params) (*rend.Root, error)
	AllowUpload(lv.Socket, params.Params) (any, error)
	Progress(lv.Socket, params.Params) (*rend.Root, error)
	DestroyCIDs([]int) error
}
