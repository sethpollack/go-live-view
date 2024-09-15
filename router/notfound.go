package router

import (
	"github.com/sethpollack/go-live-view/html"
	"github.com/sethpollack/go-live-view/rend"
	"github.com/sethpollack/go-live-view/std"
)

type notFound struct{}

func (n *notFound) Render(_ rend.Node) (rend.Node, error) {
	return html.Div(
		std.Text("404 Not Found"),
	), nil
}
