# Go LiveView

Go backend library for the Phoenix LiveView JS client. Enables rich, real-time user experiences with server-rendered HTML.


> This is still very much a work in progress. The API is not stable and is subject to change.


Basic Example:

```go
type Live struct {
	lv.Base
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

```

See the [examples](https://github.com/sethpollack/go-live-view/tree/main/examples) directory for full examples.
