package rend

import (
	"strings"

	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

type Node interface {
	Render(bool, *Root, *Rend, *strings.Builder) error
}

func RenderString(n Node) string {
	b := &strings.Builder{}
	root := NewRoot()

	render(false, root, root.Rend, b, n)

	return b.String()
}

func RenderTree(n Node) *Root {
	root := NewRoot()

	b := &strings.Builder{}

	render(true, root, root.Rend, b, n)

	root.Rend.AddStatic(b.String())

	return root
}

func RenderJSONTree(root *Root) string {
	res, err := json.Marshal(root, json.DefaultOptionsV2())
	if err != nil {
		panic(err)
	}
	(*jsontext.Value)(&res).Indent("", "\t")

	return string(res)
}

func RenderJSON(n Node) string {
	return RenderJSONTree(
		RenderTree(n),
	)
}

func Render(root *Root, n Node) *Rend {
	b := &strings.Builder{}
	rend := &Rend{}

	render(true, root, rend, b, n)

	rend.AddStatic(b.String())

	return rend
}

func render(
	diff bool,
	root *Root,
	t *Rend,
	b *strings.Builder,
	n Node,
) {
	n.Render(diff, root, t, b)
}
