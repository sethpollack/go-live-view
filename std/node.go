package std

import (
	"go-live-view/rend"
	"strings"
)

type dNode struct {
	node rend.Node
}

func DynamicNode(root rend.Node) *dNode {
	return &dNode{root}
}

func (c *dNode) Render(diff bool, root *rend.Root, t *rend.Rend, b *strings.Builder) error {
	if c.node == nil {
		return nil
	}

	if diff {
		t.AddDynamic(rend.Render(root, c.node))
		t.AddStatic(b.String())
		b.Reset()

		return nil
	}

	return c.node.Render(diff, root, t, b)
}
