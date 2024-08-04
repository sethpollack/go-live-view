package broadcast

import (
	"go-live-view/html"
	lv "go-live-view/liveview"
	"go-live-view/params"
	"go-live-view/rend"
	"go-live-view/std"
	"time"
)

type Live struct {
	lv.Base
	Time time.Time
}

func (l *Live) Mount(s lv.Socket, _ params.Params) error {
	if s != nil {
		go func() {
			time.Sleep(1 * time.Second)
			s.PushSelf("update", nil)
		}()
	}

	l.Time = time.Now()

	return nil
}

func (l *Live) Event(s lv.Socket, event string, _ params.Params) error {
	if event == "update" {
		go func() {
			time.Sleep(1 * time.Second)
			s.PushSelf("update", nil)
		}()

		l.Time = time.Now()
	}

	return nil
}

func (l *Live) Render(_ rend.Node) (rend.Node, error) {
	time := l.Time.Format("3:4:5 pm")
	return html.Div(
		html.H1(
			std.Text(&time),
		),
	), nil
}
