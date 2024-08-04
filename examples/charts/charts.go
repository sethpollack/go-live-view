package charts

import (
	"go-live-view/html"
	lv "go-live-view/liveview"
	"go-live-view/params"
	"go-live-view/rend"
	"math/rand"
	"time"

	"github.com/go-json-experiment/json"
)

type Live struct {
	lv.Base
	Options Options
}

type Chart struct {
	Type   string `json:"type"`
	Height int    `json:"height"`
}

type Series struct {
	Name string `json:"name"`
	Data []int  `json:"data"`
}

type Xaxis struct {
	Categories []int `json:"categories"`
}

type Options struct {
	Chart  Chart    `json:"chart"`
	Series []Series `json:"series"`
	Xaxis  Xaxis    `json:"xaxis"`
}

func NewOptions() Options {
	return Options{
		Chart: Chart{
			Type:   "line",
			Height: 350,
		},
		Series: []Series{
			{
				Name: "sales",
				Data: rand.Perm(9),
			},
		},
		Xaxis: Xaxis{
			Categories: []int{1991, 1992, 1993, 1994, 1995, 1996, 1997, 1998, 1999},
		},
	}
}

func (l *Live) Mount(s lv.Socket, _ params.Params) error {
	l.Options = NewOptions()
	if s != nil {
		go func() {
			time.Sleep(5 * time.Second)
			s.PushSelf("update-chart", nil)
		}()
	}
	return nil
}

func (l *Live) Event(s lv.Socket, event string, _ params.Params) error {
	if event == "update-chart" {
		go func() {
			time.Sleep(5 * time.Second)
			s.PushSelf("update-chart", nil)
		}()

		l.Options = NewOptions()
	}

	return nil
}

func (l *Live) Render(_ rend.Node) (rend.Node, error) {
	b, _ := json.Marshal(l.Options, json.DefaultOptionsV2())
	options := string(b)

	return html.Div(
		html.IdAttr("chart"),
		html.Attr("phx-hook", "Chart"),
		html.Attr("data-options", &options),
	), nil
}
