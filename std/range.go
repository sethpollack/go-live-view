package std

import (
	"go-live-view/rend"
	"strconv"
	"strings"
)

type mapRange[T any] struct {
	arr []T
	f   func(T) rend.Node
}

func Range[T any](arr []T, f func(T) rend.Node) rend.Node {
	return &mapRange[T]{
		arr: arr,
		f:   f,
	}
}

func (c *mapRange[T]) Render(diff bool, root *rend.Root, t *rend.Rend, b *strings.Builder) error {
	if len(c.arr) <= 0 {
		return nil
	}
	if !diff {
		for _, d := range c.arr {
			err := c.f(d).Render(diff, root, t, b)
			if err != nil {
				return err
			}
		}

		return nil
	}

	rends := []*rend.Rend{}

	for _, d := range c.arr {
		rend := rend.Render(root, c.f(d))
		rends = append(rends, rend)
	}

	staticsMatch := true
	for i := 1; i < len(rends); i++ {
		if !compareStatics(rends[i].Static, rends[i-1].Static) {
			staticsMatch = false
			break
		}
	}

	if staticsMatch {
		t.AddDynamic(&rend.Comprehension{
			Static:      rends[0].Static,
			Fingerprint: rends[0].Fingerprint,
			Dynamics:    copyDynamics(rends),
		})
		t.AddStatic(b.String())
		b.Reset()
	} else {
		for _, r := range rends {
			t.AddDynamic(r)
			t.AddStatic(b.String())
			b.Reset()
		}
	}

	return nil
}

func copyDynamics(d []*rend.Rend) [][]any {
	copy := [][]any{}

	for _, m := range d {
		copy = append(copy, copyDynamic(m.Dynamic))
	}

	return copy
}

func copyDynamic(m map[string]any) []any {
	copy := []any{}

	for i := 0; i < len(m); i++ {
		copy = append(copy, m[strconv.Itoa(i)])
	}

	return copy
}

func compareStatics(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i, s := range a {
		if s != b[i] {
			return false
		}
	}

	return true
}
