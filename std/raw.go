package std

import (
	"strings"

	"github.com/sethpollack/go-live-view/rend"
)

type raw struct {
	data string
}

func Raw(s string) rend.Node {
	return &raw{s}
}

func (raw *raw) Render(diff bool, root *rend.Root, t *rend.Rend, b *strings.Builder) error {
	_, err := b.Write([]byte(raw.data))
	return err
}
