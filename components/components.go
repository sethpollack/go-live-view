package components

import (
	"fmt"
	"strings"

	"github.com/sethpollack/go-live-view/html"
	"github.com/sethpollack/go-live-view/rend"
	"github.com/sethpollack/go-live-view/std"
	"github.com/sethpollack/go-live-view/uploads"
)

func Unpkg(pkg, version string) rend.Node {
	return html.Script(
		html.Attrs(
			html.ScriptDeferAttr("true"),
			html.ScriptTypeAttr("text/javascript"),
			html.ScriptSrcAttr(
				fmt.Sprintf("https://unpkg.com/%s@%s", pkg, version),
			),
		),
	)
}

func RootLayout(children ...rend.Node) rend.Node {
	return html.Html(
		html.Head(
			Unpkg("phoenix", "1.7.10"),
			Unpkg("phoenix_live_view", "0.20.3"),
			Unpkg("topbar", "2.0.2"),
			Unpkg("apexcharts", "3.26.0"),
		),
		html.Body(
			html.Div(
				children...,
			),
			html.Script(
				html.Attrs(
					html.ScriptDeferAttr("true"),
					html.Attr("phx-track-static"),
					html.ScriptTypeAttr("text/javascript"),
					html.ScriptSrcAttr("/assets/app.js"),
				),
			),
		),
	)
}

func UploadInput(u *uploads.Config, children ...rend.Node) rend.Node {
	return std.Component(
		html.Input(
			html.Attr("id", u.Ref),
			html.Attr("type", "file"),
			html.Attr("name", u.Name),
			html.Attr("accept", strings.Join(u.Accept, ",")),
			html.Attr("data-phx-hook", "Phoenix.LiveFileUpload"),
			html.Attr("data-phx-update", "ignore"),
			html.Attr("data-phx-upload-ref", u.Ref),
			html.Attr("data-phx-active-refs", u.ActiveRefs()),
			html.Attr("data-phx-done-refs", u.DoneRefs()),
			html.Attr("data-phx-preflighted-refs", u.PreflightRefs()),
			html.Attrs(
				std.TernaryNode(
					u.MaxEntries > 1,
					html.Attr("multiple"),
					nil,
				),
				std.TernaryNode(
					u.AutoUpload,
					html.Attr("data-phx-auto-upload"),
					nil,
				),
				std.Group(children...),
			),
		))
}
