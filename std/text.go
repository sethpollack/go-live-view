package std

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/sethpollack/go-live-view/rend"
)

type text struct {
	text    string
	dynamic bool
}

func Text(s any) rend.Node {
	switch val := s.(type) {
	case *string:
		return &text{text: *val, dynamic: true}
	case *int:
		return &text{text: strconv.Itoa(*val), dynamic: true}
	case string:
		return &text{text: val}
	case int:
		return &text{text: strconv.Itoa(val)}
	default:
		panic(fmt.Sprintf("invalid type %T for Text", s))
	}
}

func Textf(format string, a ...any) rend.Node {
	t := &text{}

	args := make([]interface{}, 0, len(a))

	for _, arg := range a {
		switch s := arg.(type) {
		case *string:
			t.dynamic = true
			args = append(args, *s)
		case *int:
			t.dynamic = true
			args = append(args, *s)
		default:
			args = append(args, s)
		}
	}

	t.text = fmt.Sprintf(format, args...)

	return t
}

func (text *text) Render(diff bool, root *rend.Root, t *rend.Rend, b *strings.Builder) error {
	if diff {
		if text.dynamic {
			t.AddDynamic(text.text)
			t.AddStatic(b.String())
			b.Reset()
			return nil
		}
	}

	_, err := b.Write([]byte(text.text))
	return err
}
