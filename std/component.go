package std

import (
	"go-live-view/rend"
	"strings"
)

type component struct {
	node rend.Node
}

func Component(root rend.Node) *component {
	return &component{root}
}

func (c *component) Render(diff bool, root *rend.Root, t *rend.Rend, b *strings.Builder) error {
	if c.node == nil {
		return nil
	}

	if diff {
		t.AddComponent(root, rend.Render(root, c.node))
		t.AddStatic(b.String())
		b.Reset()

		return nil
	}

	return c.node.Render(diff, root, t, b)
}
