package std

import (
	"go-live-view/rend"
	"strings"
)

type noop struct{}

func Noop() rend.Node {
	return &noop{}
}

func (n *noop) Render(diff bool, root *rend.Root, t *rend.Rend, b *strings.Builder) error {
	return nil
}
