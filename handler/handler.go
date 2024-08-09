package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/sethpollack/go-live-view/channel"
	"github.com/sethpollack/go-live-view/internal/ws"
	"github.com/sethpollack/go-live-view/lifecycle"
	lv "github.com/sethpollack/go-live-view/liveview"

	"github.com/rs/xid"
)

type handler struct {
	ctx           context.Context
	setupRoutes   func() lifecycle.Router
	setupChannels func() map[string]func() channel.Channel
}

func NewHandler(ctx context.Context,
	setupRoutes func() lifecycle.Router,
	setupChannels func() map[string]func() channel.Channel,
) *handler {
	return &handler{
		ctx:           ctx,
		setupRoutes:   setupRoutes,
		setupChannels: setupChannels,
	}
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	hub := channel.NewHub()
	go hub.Listen(h.ctx)

	ws.NewWSHandler(
		func(w http.ResponseWriter, r *http.Request) {
			rt := h.setupRoutes()

			lc := lifecycle.NewLifecycle(rt)

			resp, err := lc.StaticRender(xid.New().String(), r.URL.Path)

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprintf("Error: %s", err.Error())))
				return
			}

			w.Write([]byte(resp))
		},
		func(t channel.Transport) {
			rt := h.setupRoutes()

			lc := lifecycle.NewLifecycle(rt)

			server := channel.NewServer(t, hub)

			hub.Add(server)
			defer hub.Remove(server)

			server.Route("lv:*", lv.NewLVChannel(lc))
			server.Route("lvu:*", lv.NewLVUChannel(lc))

			if h.setupChannels != nil {
				for topic, factory := range h.setupChannels() {
					server.Route(topic, factory)
				}
			}

			server.Listen(h.ctx)
		},
	).Serve(w, r)
}
