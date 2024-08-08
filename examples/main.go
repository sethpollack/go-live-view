package main

import (
	"context"
	comp "go-live-view/components"
	"go-live-view/examples/async"
	"go-live-view/examples/broadcast"
	"go-live-view/examples/charts"
	"go-live-view/examples/comprehension"
	"go-live-view/examples/counter"
	"go-live-view/examples/index"
	"go-live-view/examples/nested"
	"go-live-view/examples/scroll"
	"go-live-view/examples/ssnav"
	"go-live-view/examples/stream"
	"go-live-view/examples/uploads"
	"go-live-view/handler"
	"go-live-view/lifecycle"
	"go-live-view/router"
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
	rt := router.NewRouter(
		router.WithLayout(comp.Layout),
	)

	root := rt.Group("/", &index.Live{
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
	})

	root.Handle("/counter", &counter.Live{})
	root.Handle("/uploads", uploads.New())
	root.Handle("/chart", &charts.Live{})
	root.Handle("/async", &async.Live{})
	root.Handle("/broadcast", broadcast.New())
	root.Handle("/comprehension", &comprehension.Live{})
	root.Handle("/stream", &stream.Live{})
	root.Handle("/scroll", &scroll.Live{})

	nest := root.Group("/nested", &nested.Live{})
	nest.Handle("/:id", &nested.ShowLive{})
	nest.Handle("/:id/edit", &nested.EditLive{})

	snav := root.Group("/ssnav", &ssnav.Live{})
	snav.Handle("/:id", &ssnav.ShowLive{})
	snav.Handle("/:id/edit", &ssnav.EditLive{})

	return rt
}

func main() {
	ctx := context.Background()

	mux := http.NewServeMux()

	mux.Handle("/assets/app.js", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(appJS))
	}))

	mux.Handle("/favicon.ico", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(""))
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

	log.Println("server listening on", srv.Addr)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	srv.Shutdown(ctx)

	log.Println("shutting down")
	os.Exit(0)
}