package std

import (
	"strings"

	"github.com/sethpollack/go-live-view/rend"
)

type goEmbed struct {
	cb func() rend.Node
}

func GoEmbed(cb func() rend.Node) rend.Node {
	return &goEmbed{cb}
}

func (v *goEmbed) Render(diff bool, root *rend.Root, t *rend.Rend, b *strings.Builder) error {
	node := v.cb()

	if node == nil {
		return nil
	}

	if diff {
		t.AddDynamic(rend.Render(root, node))
		t.AddStatic(b.String())
		b.Reset()

		return nil
	}

	return node.Render(diff, root, t, b)
}
