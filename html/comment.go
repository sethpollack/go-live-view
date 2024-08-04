package html

import (
	"go-live-view/rend"
	"strings"
)

type comment struct {
	comment string
}

func Comment(s string) rend.Node {
	return &comment{s}
}

func (c *comment) Render(diff bool, root *rend.Root, t *rend.Rend, b *strings.Builder) error {
	b.WriteString("<!--")
	b.WriteString(c.comment)
	b.WriteString("-->")
	return nil
}
