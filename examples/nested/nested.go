package nested

import (
	"github.com/sethpollack/go-live-view/html"
	lv "github.com/sethpollack/go-live-view/liveview"
	"github.com/sethpollack/go-live-view/params"
	"github.com/sethpollack/go-live-view/rend"
	"github.com/sethpollack/go-live-view/std"
)

type Live struct {
}

func (u *Live) Render(child rend.Node) (rend.Node, error) {
	return html.Div(
		html.H1(
			std.Text("Nested"),
		),
		html.Button(
			html.A(
				std.Text("Show"),
				html.AHrefAttr("/nested/1"),
				html.DataAttr("phx-link", "patch"),
				html.DataAttr("phx-link-state", "push"),
			),
		),
		html.Button(
			html.A(
				std.Text("Edit"),
				html.AHrefAttr("/nested/1/edit"),
				html.DataAttr("phx-link", "patch"),
				html.DataAttr("phx-link-state", "push"),
			),
		),
		child,
	), nil
}

type ShowLive struct {
	params params.Params
}

func (l *ShowLive) Params(s lv.Socket, p params.Params) error {
	l.params = p
	return nil
}

func (l *ShowLive) Render(_ rend.Node) (rend.Node, error) {
	id := l.params.String("id")
	return html.Div(
		html.H1(
			std.Textf("Show %s", &id),
		),
	), nil
}

type EditLive struct {
	params params.Params
}

func (l *EditLive) Params(s lv.Socket, p params.Params) error {
	l.params = p
	return nil
}

func (l *EditLive) Render(_ rend.Node) (rend.Node, error) {
	id := l.params.String("id")
	return html.Div(
		html.H1(
			std.Textf("Edit %s", &id),
		),
	), nil
}
