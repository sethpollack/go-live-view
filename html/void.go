package html

import (
	"strings"

	"github.com/sethpollack/go-live-view/rend"
)

type void struct {
	tag   string
	attrs []rend.Node
}

func Void(tag string, children ...rend.Node) *void {
	v := &void{
		tag: tag,
	}

	for _, attr := range children {
		if attr == nil {
			continue
		}

		switch attr.(type) {
		case *attrs, *attribute:
			v.attrs = append(v.attrs, attr)
		}
	}

	return v
}

func (v *void) Render(diff bool, root *rend.Root, t *rend.Rend, b *strings.Builder) error {
	_, err := b.Write([]byte("<" + v.tag))
	if err != nil {
		return err
	}

	for _, child := range v.attrs {
		err := child.Render(diff, root, t, b)
		if err != nil {
			return err
		}
	}

	_, err = b.Write([]byte("/>"))
	if err != nil {
		return err
	}

	return nil
}
