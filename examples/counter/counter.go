package counter

import (
	"github.com/sethpollack/go-live-view/html"
	lv "github.com/sethpollack/go-live-view/liveview"
	"github.com/sethpollack/go-live-view/params"
	"github.com/sethpollack/go-live-view/rend"
	"github.com/sethpollack/go-live-view/std"
)

type Live struct {
	Count int
}

func (l *Live) Event(s lv.Socket, event string, _ params.Params) error {

	if event == "inc" {
		l.Count++
	}

	if event == "dec" {
		l.Count--
	}

	return nil
}

func (l *Live) Render(_ rend.Node) (rend.Node, error) {
	return html.Div(
		html.H1(
			std.Text(&l.Count),
		),
		html.Button(
			std.Text("inc"),
			html.Attr("phx-click", "inc"),
		),
		html.Button(
			std.Text("dec"),
			html.Attr("phx-click", "dec"),
		),
	), nil
}
