package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/rs/xid"
	"github.com/sethpollack/go-live-view/channel"
	"github.com/sethpollack/go-live-view/channel/transport/longpoll"
	"github.com/sethpollack/go-live-view/channel/transport/websocket"
	"github.com/sethpollack/go-live-view/lifecycle"
	lv "github.com/sethpollack/go-live-view/liveview"
)

type handlerOption func(*handler)

func WithChannels(channels map[string]func() channel.Channel) handlerOption {
	return func(h *handler) {
		h.channels = channels
	}
}

func WithTransports(transport channel.Transport) handlerOption {
	return func(h *handler) {
		h.transports = append(h.transports, transport)
	}
}

type handler struct {
	ctx         context.Context
	setupRoutes func() lifecycle.Router
	channels    map[string]func() channel.Channel
	channelHub  *channel.Hub
	transports  []channel.Transport
}

func NewHandler(ctx context.Context, setupRoutes func() lifecycle.Router, opts ...handlerOption) *handler {
	h := &handler{
		ctx:         ctx,
		setupRoutes: setupRoutes,
		channelHub:  channel.NewHub(),
		transports: []channel.Transport{
			websocket.New("/live/websocket"),
			longpoll.New("/live/longpoll"),
		},
	}

	go h.channelHub.Listen(h.ctx)

	for _, opt := range opts {
		opt(h)
	}

	return h
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, transport := range h.transports {
		if transport.Path() == r.URL.Path {
			transport.Serve(h.handle, w, r)
			return
		}
	}

	resp, err := lifecycle.NewLifecycle(h.setupRoutes()).
		StaticRender(xid.New().String(), r.URL.Path)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Error: %s", err.Error())))
		return
	}

	w.Write([]byte(resp))
}

func (h *handler) handle(t channel.Conn) {
	server := channel.NewServer(t, h.channelHub)
	h.channelHub.Add(server)
	defer h.channelHub.Remove(server)

	rt := h.setupRoutes()
	lc := lifecycle.NewLifecycle(rt)

	server.Route("lv:*", lv.NewLVChannel(lc))
	server.Route("lvu:*", lv.NewLVUChannel(lc))

	for topic, factory := range h.channels {
		server.Route(topic, factory)
	}

	server.Listen(h.ctx)
}
