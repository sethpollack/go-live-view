package liveview

import (
	"errors"

	"github.com/sethpollack/go-live-view/params"
	"github.com/sethpollack/go-live-view/rend"
)

var NotFoundError = errors.New("route not found")

type Route interface {
	GetView() View
	GetParams() params.Params
}

type Router interface {
	GetRoute(string) (Route, error)
	Routable(Route, Route) bool
	GetLayout() func(...rend.Node) rend.Node
}
