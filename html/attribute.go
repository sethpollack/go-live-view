package html

import (
	"fmt"
	"go-live-view/rend"
	"sort"
	"strings"
)

type attribute struct {
	tag     string
	value   string
	dynamic bool
}

func Attr(tagName string, value ...any) rend.Node {
	return tag(tagName, value...)
}

func (attr *attribute) Render(diff bool, root *rend.Root, t *rend.Rend, b *strings.Builder) error {
	if !diff || !attr.dynamic {
		if attr.value == "" {
			_, err := b.Write([]byte(fmt.Sprintf(" %s", attr.tag)))
			return err
		}
		_, err := b.Write([]byte(fmt.Sprintf(" %s=\"%s\"", attr.tag, attr.value)))
		return err
	}

	_, err := b.Write([]byte(fmt.Sprintf(" %s=", attr.tag)))
	if err != nil {
		return err
	}

	if attr.value == "" {
		t.AddDynamic("\"\"")
	} else {
		t.AddDynamic(attr.value)
	}

	t.AddStatic(b.String())
	b.Reset()

	return nil
}

func tag(tag string, value ...any) rend.Node {
	if len(value) == 0 {
		return &attribute{
			tag: tag,
		}
	}

	switch val := value[0].(type) {
	case *string:
		return &attribute{
			tag:     tag,
			value:   *val,
			dynamic: true,
		}
	case *int:
		return &attribute{
			tag:     tag,
			value:   fmt.Sprintf("%d", *val),
			dynamic: true,
		}
	case *bool:
		return &attribute{
			tag:     tag,
			value:   fmt.Sprintf("%t", *val),
			dynamic: true,
		}
	case bool:
		return &attribute{
			tag:   tag,
			value: fmt.Sprintf("%t", val),
		}
	case int, float32, float64:
		return &attribute{
			tag:   tag,
			value: fmt.Sprintf("%d", val),
		}
	case string:
		return &attribute{
			tag:   tag,
			value: val,
		}
	default:
		panic(fmt.Sprintf("invalid type %T for Attr", value[0]))
	}
}

type attrs struct {
	Attrs []rend.Node
}

func Attrs(children ...rend.Node) rend.Node {
	return &attrs{Attrs: children}
}

func Values(prefix string, values map[string]any) rend.Node {
	attrs := []rend.Node{}
	keys := make([]string, 0, len(values))

	for k := range values {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		v := values[k]
		attrs = append(attrs, Attr(prefix+k, v))
	}

	return Attrs(attrs...)
}

func (g *attrs) Render(diff bool, root *rend.Root, t *rend.Rend, b *strings.Builder) error {
	for _, child := range g.Attrs {
		err := child.Render(diff, root, t, b)
		if err != nil {
			return err
		}
	}
	return nil
}
