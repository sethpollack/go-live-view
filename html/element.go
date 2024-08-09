package html

import (
	"strings"

	"github.com/sethpollack/go-live-view/rend"
)

type element struct {
	tag      string
	children []rend.Node
	attrs    []rend.Node
}

func Element(tag string, children ...rend.Node) rend.Node {
	e := &element{
		tag: tag,
	}

	for _, child := range children {
		switch child.(type) {
		case *attrs, *attribute:
			e.attrs = append(e.attrs, child)
		default:
			e.children = append(e.children, child)
		}
	}

	return e
}

func (el *element) Render(diff bool, root *rend.Root, t *rend.Rend, b *strings.Builder) error {
	_, err := b.Write([]byte("<" + el.tag))
	if err != nil {
		return err
	}

	for _, child := range el.attrs {
		if child == nil {
			continue
		}

		err := child.Render(diff, root, t, b)
		if err != nil {
			return err
		}
	}

	_, err = b.Write([]byte(">"))
	if err != nil {
		return err
	}

	for _, child := range el.children {
		if child == nil {
			continue
		}

		err := child.Render(diff, root, t, b)
		if err != nil {
			return err
		}
	}

	_, err = b.Write([]byte("</" + el.tag + ">"))
	if err != nil {
		return err
	}

	return nil
}
