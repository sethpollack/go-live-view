package ssnav

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

func (u *Live) Event(s lv.Socket, e string, p params.Params) error {
	if e == "navigate" {
		return s.PushPatch(p.Map("value").String("href"), false)
	}
	return nil
}

func (u *Live) Render(child rend.Node) (rend.Node, error) {
	return html.Div(
		html.H1(
			std.Text("Server Navigation"),
		),
		html.Button(
			html.A(
				std.Text("Show"),
				html.Attr("phx-click", "navigate"),
				html.Attr("phx-value-href", "/ssnav/1"),
			),
		),
		html.Button(
			html.A(
				std.Text("Edit"),
				html.Attr("phx-click", "navigate"),
				html.Attr("phx-value-href", "/ssnav/1/edit"),
			),
		),
		child,
	), nil
}

type ShowLive struct {
	lv.Base
	id string
}

func (l *ShowLive) Params(s lv.Socket, p params.Params) error {
	l.id = p.String("id")
	return nil
}

func (l *ShowLive) Render(_ rend.Node) (rend.Node, error) {
	return html.Div(
		html.H1(
			std.Textf("Show %s", &l.id),
		),
	), nil
}

type EditLive struct {
	lv.Base
	id string
}

func (l *EditLive) Params(s lv.Socket, p params.Params) error {
	l.id = p.String("id")
	return nil
}

func (l *EditLive) Render(_ rend.Node) (rend.Node, error) {
	return html.Div(
		html.H1(
			std.Textf("Edit %s", &l.id),
		),
	), nil
}
