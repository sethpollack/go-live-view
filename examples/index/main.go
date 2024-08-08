package index

import (
	"go-live-view/html"
	lv "go-live-view/liveview"
	"go-live-view/rend"
	"go-live-view/std"
)

type Live struct {
	lv.Base
	Links []string
}

func (i *Live) Render(child rend.Node) (rend.Node, error) {
	return html.Div(
		html.Ol(
			std.Range(i.Links, func(link string) rend.Node {
				return html.Li(
					html.A(
						html.AHrefAttr(&link),
						std.Text(&link),
						html.DataAttr("phx-link", "patch"),
						html.DataAttr("phx-link-state", "push"),
					),
				)
			}),
		),
		html.Div(
			std.DynamicNode(child),
		),
	), nil
}
