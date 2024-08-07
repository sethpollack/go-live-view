package main

import (
	"context"
	comp "go-live-view/components"
	"go-live-view/examples/async"
	"go-live-view/examples/broadcast"
	"go-live-view/examples/charts"
	"go-live-view/examples/comprehension"
	"go-live-view/examples/counter"
	"go-live-view/examples/nested"
	"go-live-view/examples/scroll"
	"go-live-view/examples/ssnav"
	"go-live-view/examples/stream"
	"go-live-view/examples/uploads"
	"go-live-view/handler"
	"go-live-view/html"
	"go-live-view/lifecycle"
	lv "go-live-view/liveview"
	"go-live-view/rend"
	router "go-live-view/routers/nested"
	"go-live-view/std"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

const appJS = `
(() => {
	window.addEventListener("phx:page-loading-start", info => topbar.show())
	window.addEventListener("phx:page-loading-stop", info => topbar.hide())

	let Hooks = {}
	Hooks.Chart = {
		mounted() {
			const options = JSON.parse(this.el.dataset.options)
			window.chart = new ApexCharts(this.el, options);
			window.chart.render();
		},
		updated() {
			const options = JSON.parse(this.el.dataset.options)
			window.chart.updateSeries(options.series)
		}
	}

	const lv = new LiveView.LiveSocket("/live", Phoenix.Socket, {hooks: Hooks});
	lv.connect();
})();
`

func setupRoutes() lifecycle.Router {
	rt := router.NewRouter()

	rt.HandleLive(&router.Routes{
		Route: &router.Route{
			Path: "/",
			View: &IndexLive{
				Links: []string{
					"/async",
					"/broadcast",
					"/chart",
					"/comprehension",
					"/counter",
					"/nested",
					"/ssnav",
					"/scroll",
					"/stream",
					"/uploads",
				},
			},
			Layout: comp.Layout,
		},
		Children: []*router.Routes{
			{
				Route: &router.Route{
					Path: "/counter",
					View: &counter.Live{},
				},
			},
			{
				Route: &router.Route{
					Path: "/uploads",
					View: uploads.New(),
				},
			},
			{
				Route: &router.Route{
					Path: "/chart",
					View: &charts.Live{},
				},
			},
			{
				Route: &router.Route{
					Path: "/async",
					View: &async.Live{},
				},
			},
			{
				Route: &router.Route{
					Path: "/broadcast",
					View: broadcast.New(),
				},
			},
			{
				Route: &router.Route{
					Path: "/comprehension",
					View: &comprehension.Live{},
				},
			},
			{
				Route: &router.Route{
					Path: "/stream",
					View: &stream.Live{},
				},
			},
			{
				Route: &router.Route{
					Path: "/scroll",
					View: &scroll.Live{},
				},
			},
			{
				Route: &router.Route{
					Path: "/nested",
					View: &nested.Live{},
				},
				Children: []*router.Routes{
					{
						Route: &router.Route{
							Path: "/:id",
							View: &nested.ShowLive{},
						},
					},
					{
						Route: &router.Route{
							Path: "/:id/edit",
							View: &nested.EditLive{},
						},
					},
				},
			},
			{
				Route: &router.Route{
					Path: "/ssnav",
					View: &ssnav.Live{},
				},
				Children: []*router.Routes{
					{
						Route: &router.Route{
							Path: "/:id",
							View: &ssnav.ShowLive{},
						},
					},
					{
						Route: &router.Route{
							Path: "/:id/edit",
							View: &ssnav.EditLive{},
						},
					},
				},
			},
		},
	})

	return rt
}

func main() {
	ctx := context.Background()

	mux := http.NewServeMux()

	mux.Handle("/assets/app.js", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(appJS))
	}))

	mux.Handle("/", handler.NewHandler(ctx, setupRoutes, nil))

	srv := &http.Server{
		Addr: "0.0.0.0:8080",

		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      mux,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	log.Println("server listening on :8080")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	srv.Shutdown(ctx)

	log.Println("shutting down")
	os.Exit(0)
}

type IndexLive struct {
	lv.Base
	Links []string
}

func (i *IndexLive) Render(child rend.Node) (rend.Node, error) {
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
