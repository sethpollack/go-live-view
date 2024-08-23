package comprehension

import (
	"fmt"

	"github.com/sethpollack/go-live-view/html"
	lv "github.com/sethpollack/go-live-view/liveview"
	"github.com/sethpollack/go-live-view/params"
	"github.com/sethpollack/go-live-view/rend"
	"github.com/sethpollack/go-live-view/std"
)

type User struct {
	ID   int
	Name string
}

type Live struct {
	users []*User
}

func (l *Live) Event(s lv.Socket, event string, p params.Params) error {
	if event == "add-user" {
		l.users = append(l.users,
			&User{
				ID:   len(l.users) + 1,
				Name: fmt.Sprintf("User %d", len(l.users)+1),
			},
		)
	}

	if event == "delete-user" {
		id := p.Map("value").Int("id")

		for i, u := range l.users {
			if u.ID == id {
				l.users = append(l.users[:i], l.users[i+1:]...)
				break
			}
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
				std.Range[*User](l.users, func(u *User) rend.Node {
					return html.Tr(
						html.Td(
							std.Text(&u.Name),
							html.Button(
								html.Attrs(
									html.Attr("phx-click", "delete-user"),
									html.Attr("phx-value-id", &u.ID),
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
