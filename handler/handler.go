package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/sethpollack/go-live-view/channel"
	"github.com/sethpollack/go-live-view/channel/transport/longpoll"
	"github.com/sethpollack/go-live-view/channel/transport/websocket"
	"github.com/sethpollack/go-live-view/internal/lvchan"
	"github.com/sethpollack/go-live-view/internal/lvuchan"
	lv "github.com/sethpollack/go-live-view/liveview"
)

type handlerOption func(*handler)

type handler struct {
	ctx           context.Context
	setupRoutes   func() lv.Router
	channels      map[string]func() channel.Channel
	channelHub    *channel.Hub
	transports    []channel.Transport
	tokenizer     tokenizer
	sessionGetter sessionGetter
}

func NewHandler(ctx context.Context, setupRoutes func() lv.Router, opts ...handlerOption) *handler {
	h := &handler{
		ctx:         ctx,
		setupRoutes: setupRoutes,
		channelHub:  channel.NewHub(),
		channels:    make(map[string]func() channel.Channel),
		transports: []channel.Transport{
			websocket.New("/live/websocket"),
			longpoll.New("/live/longpoll"),
		},
		tokenizer:     &defaultTokenizer{},
		sessionGetter: &defaultSessionGetter{},
	}

	go h.channelHub.Listen(h.ctx)

	for _, opt := range opts {
		opt(h)
	}

	return h
}

func WithChannel(pattern string, f func() channel.Channel) handlerOption {
	return func(h *handler) {
		h.channels[pattern] = f
	}
}

func WithTransport(transport channel.Transport) handlerOption {
	return func(h *handler) {
		h.transports = append(h.transports, transport)
	}
}

func WithTokenizer(tokenizer tokenizer) handlerOption {
	return func(h *handler) {
		h.tokenizer = tokenizer
	}
}

func WithSessionGetter(getter sessionGetter) handlerOption {
	return func(h *handler) {
		h.sessionGetter = getter
	}
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, transport := range h.transports {
		if transport.Path() == r.URL.Path {
			transport.Serve(h.handle, w, r)
			return
		}
	}

	resp, err := lv.NewLifecycle(
		h.setupRoutes(), h.tokenizer, h.sessionGetter,
	).StaticRender(w, r)
	if err != nil {
		switch err.(type) {
		case lv.HttpError:
			httpErr := err.(lv.HttpError)
			w.WriteHeader(httpErr.Code())
			w.Write([]byte(httpErr.Error()))
			return
		default:
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Error: %s", err.Error())))
			return
		}
	}

	w.Write([]byte(resp))
}

func (h *handler) handle(t channel.Conn) {
	server := channel.NewServer(t, h.channelHub)
	h.channelHub.Add(server)
	defer h.channelHub.Remove(server)

	rt := h.setupRoutes()
	lc := lv.NewLifecycle(rt, h.tokenizer, h.sessionGetter)

	server.Route("lv:*", lvchan.New(lc))
	server.Route("lvu:*", lvuchan.New(lc))

	for topic, factory := range h.channels {
		server.Route(topic, factory)
	}

	server.Listen(h.ctx)
}
