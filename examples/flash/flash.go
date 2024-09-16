package flash

import (
	"github.com/sethpollack/go-live-view/html"
	lv "github.com/sethpollack/go-live-view/liveview"
	"github.com/sethpollack/go-live-view/params"
	"github.com/sethpollack/go-live-view/rend"
	"github.com/sethpollack/go-live-view/std"
)

type Live struct {
	flashes []string
}

func (l *Live) Event(s lv.Socket, event string, p params.Params) error {
	switch event {
	case "patch":
		s.PushPatch("/", lv.WithFlash("info", "from patch!"))
	case "navigate":
		s.PushNavigate("/", lv.WithFlash("info", "from navigate!"))
	case "redirect":
		s.Redirect("/", lv.WithFlash("info", "from redirect!"))
	case "no-flash-navigate":
		s.PushNavigate("/")
	}

	return nil
}

func (l *Live) Render(_ rend.Node) (rend.Node, error) {
	return html.Div(
		html.Button(
			std.Text("patch"),
			html.Attr("phx-click", "patch"),
		),
		html.Button(
			std.Text("navigate"),
			html.Attr("phx-click", "navigate"),
		),
		html.Button(
			std.Text("redirect"),
			html.Attr("phx-click", "redirect"),
		),
		html.Button(
			std.Text("flash"),
			html.Attr("phx-click", "lv:flash"),
			html.Attr("phx-value-key", "info"),
			html.Attr("phx-value-msg", "from event"),
		),
		html.Button(
			std.Text("no flash navigate"),
			html.Attr("phx-click", "no-flash-navigate"),
		),
	), nil
}
