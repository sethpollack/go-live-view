package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	comp "github.com/sethpollack/go-live-view/components"
	"github.com/sethpollack/go-live-view/examples/async"
	"github.com/sethpollack/go-live-view/examples/broadcast"
	"github.com/sethpollack/go-live-view/examples/charts"
	"github.com/sethpollack/go-live-view/examples/comprehension"
	"github.com/sethpollack/go-live-view/examples/counter"
	"github.com/sethpollack/go-live-view/examples/index"
	"github.com/sethpollack/go-live-view/examples/js"
	"github.com/sethpollack/go-live-view/examples/nested"
	"github.com/sethpollack/go-live-view/examples/scroll"
	"github.com/sethpollack/go-live-view/examples/ssnav"
	"github.com/sethpollack/go-live-view/examples/stream"
	"github.com/sethpollack/go-live-view/examples/uploads"
	"github.com/sethpollack/go-live-view/handler"
	lv "github.com/sethpollack/go-live-view/liveview"
	"github.com/sethpollack/go-live-view/router"
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

	window.liveSocket = lv;
})();
`

func setupRoutes() lv.Router {
	rt := router.NewRouter(
		comp.RootLayout,
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
			"/js",
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
	root.Handle("/js", &js.Live{})

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

	mux.Handle("/", handler.NewHandler(ctx, setupRoutes))

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
