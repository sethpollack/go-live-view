package longpoll

import (
	"net/http"

	"github.com/sethpollack/go-live-view/channel"
)

var _ channel.Transport = (*lpTransport)(nil)

func New(path string) channel.Transport {
	return &lpTransport{
		path: path,
	}
}

type lpTransport struct {
	path string
}

func (l *lpTransport) Path() string {
	return l.path
}

func (t *lpTransport) Serve(handle func(channel.Conn), w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
}
