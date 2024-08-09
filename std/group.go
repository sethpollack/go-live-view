package std

import (
	"strings"

	"github.com/sethpollack/go-live-view/rend"
)

type group struct {
	Children []rend.Node
}

func Group(children ...rend.Node) rend.Node {
	return &group{Children: children}
}

func (group *group) Render(diff bool, root *rend.Root, t *rend.Rend, b *strings.Builder) error {
	for _, child := range group.Children {
		err := child.Render(diff, root, t, b)
		if err != nil {
			return err
		}
	}
	return nil
}
