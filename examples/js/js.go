package js

import (
	"github.com/sethpollack/go-live-view/html"
	"github.com/sethpollack/go-live-view/js"
	lv "github.com/sethpollack/go-live-view/liveview"
	"github.com/sethpollack/go-live-view/rend"
	"github.com/sethpollack/go-live-view/std"
)

type Live struct {
	*lv.Base
}

func hide_modal() string {
	return js.JS(
		js.Hide(&js.HideArgs{
			To:         "#modal",
			Transition: [3]string{"fade-out"},
		}),
		js.Hide(&js.HideArgs{
			To:         "#modal-content",
			Transition: [3]string{"fade-out-scale"},
		}),
	)
}

func (l *Live) Render(_ rend.Node) (rend.Node, error) {
	return html.Div(
		std.Text("Example JS Modal"),
		html.Div(
			html.IdAttr("modal"),
			html.Attr("phx-remove", hide_modal()),
			html.Div(
				html.IdAttr("modal-content"),
				html.Attr("phx-click-away", hide_modal()),
				html.Attr("phx-window-keydown", hide_modal()),
				html.Attr("phx-key", "escape"),
				std.Text("Model Content"),
				html.Button(
					html.ClassAttr("phx-modal-close"),
					html.Attr("phx-click", hide_modal()),
					std.Text("✖"),
				),
			),
		),
	), nil
}
