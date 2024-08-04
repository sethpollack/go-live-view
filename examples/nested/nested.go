package nested

import (
	"go-live-view/html"
	lv "go-live-view/liveview"
	"go-live-view/params"
	"go-live-view/rend"
	"go-live-view/std"
)

type Live struct {
	lv.Base
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
	lv.Base
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
	lv.Base
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
