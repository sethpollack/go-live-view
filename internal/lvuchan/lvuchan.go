package lvuchan

import (
	"fmt"
	"strings"

	"github.com/sethpollack/go-live-view/channel"
	"github.com/sethpollack/go-live-view/params"
)

var _ channel.Channel = &lvuChannel{}

type lvuLifecycle interface {
	Chunk(string, string, []byte, func() error) error
}

type lvuChannel struct {
	lc        lvuLifecycle
	configRef string
	ref       string
}

func New(lc lvuLifecycle) func() channel.Channel {
	return func() channel.Channel {
		return &lvuChannel{
			lc: lc,
		}
	}
}

func (l *lvuChannel) Join(s channel.Socket, p any) error {
	token := params.FromAny(p).String("token")
	splits := strings.Split(token, "-")
	if len(splits) != 2 {
		return fmt.Errorf("invalid token")
	}
	l.configRef = splits[0]
	l.ref = splits[1]

	return s.Push("", nil)
}

func (l *lvuChannel) Leave(s channel.Socket) error {
	return s.Push("", nil)
}

func (l *lvuChannel) Broadcast(s channel.Socket, event string, p any) error {
	return s.Push("", nil)
}

func (l *lvuChannel) Message(s channel.Socket, event string, p any) error {
	if event == "chunk" {
		data, ok := p.([]byte)
		if !ok {
			return fmt.Errorf("invalid chunk data")
		}

		err := l.lc.Chunk(l.configRef, l.ref, data, s.Close)
		if err != nil {
			return err
		}
	}

	return s.Push("", nil)
}
