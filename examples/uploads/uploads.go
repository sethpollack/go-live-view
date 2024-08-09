package uploads

import (
	"fmt"

	comp "github.com/sethpollack/go-live-view/components"
	"github.com/sethpollack/go-live-view/html"
	lv "github.com/sethpollack/go-live-view/liveview"
	"github.com/sethpollack/go-live-view/params"
	"github.com/sethpollack/go-live-view/rend"
	"github.com/sethpollack/go-live-view/std"
	"github.com/sethpollack/go-live-view/uploads"
)

type Live struct {
	lv.Base
	uploads *uploads.Uploads
}

func New() *Live {
	return &Live{
		uploads: uploads.New(),
	}
}

func (l *Live) Mount(s lv.Socket, p params.Params) error {
	l.uploads.AllowUpload("mydoc",
		uploads.WithAccept(".pdf"),
		uploads.WithAutoUpload(false),
		uploads.WithMaxEntries(1),
	)
	return nil
}

func (l *Live) Event(s lv.Socket, event string, p params.Params) error {
	if event == "validate" {
		l.uploads.OnValidate(p)
	}

	if event == "save" {
		l.uploads.Consume("mydoc", func(path string, entry *uploads.Entry) {
			fmt.Printf("Consuming %s", entry.Meta.Name)
		})
	}

	return nil
}

func (l *Live) Uploads() *uploads.Uploads {
	return l.uploads
}

func (l *Live) Render(_ rend.Node) (rend.Node, error) {
	return std.Component(
		html.Div(
			html.Form(
				html.Attr("id", "upload-form"),
				html.Attr("phx-submit", "save"),
				html.Attr("phx-change", "validate"),
				comp.UploadInput(l.uploads.GetByName("mydoc")),
				html.Button(
					html.Attr("type", "submit"),
					std.Text("Upload"),
				),
			),
		)), nil
}
