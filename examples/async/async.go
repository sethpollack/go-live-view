package async

import (
	"time"

	"github.com/sethpollack/go-live-view/async"
	"github.com/sethpollack/go-live-view/html"
	lv "github.com/sethpollack/go-live-view/liveview"
	"github.com/sethpollack/go-live-view/params"
	"github.com/sethpollack/go-live-view/rend"
	"github.com/sethpollack/go-live-view/std"
)

type User struct {
	Name string
}

type Live struct {
	User *async.Async[*User]
}

func (l *Live) Mount(s lv.Socket, _ params.Params) error {
	l.User = async.New(s, func() (*User, error) {
		time.Sleep(2 * time.Second)
		return &User{Name: "John"}, nil
	})

	return nil
}

func (l *Live) Render(_ rend.Node) (rend.Node, error) {
	return html.Div(
		html.H1(
			std.GoEmbed(func() rend.Node {
				switch l.User.State() {
				case async.Loading:
					loadingMessage := "Loading..."
					return std.Text(&loadingMessage)
				case async.Failed:
					err := l.User.Error().Error()
					return std.Textf("failed to load user: %s", &err)
				default:
					return std.Text(&l.User.Value().Name)
				}
			}),
		),
	), nil
}
