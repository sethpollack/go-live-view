package stream

import (
	"fmt"
	"go-live-view/html"
	"go-live-view/internal/ref"
	lv "go-live-view/liveview"
	"go-live-view/params"
	"go-live-view/rend"
	"go-live-view/std"
	"go-live-view/stream"
)

type Live struct {
	ref *ref.Ref
	lv.Base

	userStream *stream.StreamGetter
}

type User struct {
	id   int
	Name string
}

func NewUser(id int) *User {
	return &User{id: id, Name: fmt.Sprintf("User %d", id)}
}

func (l *Live) Mount(s lv.Socket, _ params.Params) error {
	l.ref = ref.New(0)

	l.userStream = stream.New("users",
		stream.IDFunc(func(item any) string {
			user := item.(*User)
			return fmt.Sprintf("user-%d", user.id)
		}),
	)

	return nil
}

func (l *Live) Event(s lv.Socket, event string, p params.Params) error {
	if event == "add-user" {
		err := l.userStream.Add(
			NewUser(int(l.ref.NextRef())),
		)
		if err != nil {
			return fmt.Errorf("adding user in event: %w", err)
		}
	}

	if event == "delete-user" {
		err := l.userStream.Delete(p.Map("value").String("id"))
		if err != nil {
			return fmt.Errorf("deleting user in event: %w", err)
		}
	}

	return nil
}

func (l *Live) Render(_ rend.Node) (rend.Node, error) {
	return html.Div(
		html.Button(
			std.Text("Add User"),
			html.Attrs(
				html.Attr("phx-click", "add-user"),
			),
		),
		html.Table(
			html.Tbody(
				html.Attrs(
					html.IdAttr("stream-users"),
					html.Attr("phx-update", "stream"),
				),
				std.Stream(l.userStream.Get(), func(item stream.Item) rend.Node {
					u := item.Item.(*User)
					return html.Tr(
						html.Attrs(
							html.IdAttr(&item.DomID),
						),
						html.Td(
							std.Text(&u.Name),
							html.Button(
								html.Attrs(
									html.Attr("phx-click", "delete-user"),
									html.Attr("phx-value-id", &item.DomID),
								),
								std.Text("Delete"),
							),
						),
					)
				}),
			),
		),
	), nil
}
