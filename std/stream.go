package std

import (
	"go-live-view/rend"
	s "go-live-view/stream"
	"strings"
)

type stream struct {
	stream *s.Stream
	f      func(s.Item) rend.Node
}

func Stream(s *s.Stream, f func(s.Item) rend.Node) *stream {
	return &stream{
		stream: s,
		f:      f,
	}
}

func (s *stream) Render(diff bool, root *rend.Root, t *rend.Rend, b *strings.Builder) error {
	if s.stream == nil {
		return nil
	}

	if len(s.stream.Deletions) == 0 && len(s.stream.Additions) == 0 && !s.stream.Reset {
		return nil
	}

	if !diff {
		for _, d := range s.stream.Additions {
			err := s.f(d).Render(diff, root, t, b)
			if err != nil {
				return err
			}
		}
		return nil
	}

	if len(s.stream.Additions) == 0 {
		t.AddDynamic(&rend.Comprehension{
			Stream: []any{
				root.NextStreamID(),
				[]any{},
				s.stream.Deletions,
				s.stream.Reset,
			},
		})
		t.AddStatic(b.String())
		b.Reset()
		return nil
	}

	rends := []*rend.Rend{}

	for _, d := range s.stream.Additions {
		rend := rend.Render(root, s.f(d))
		rends = append(rends, rend)
	}

	staticsMatch := true
	for i := 1; i < len(rends); i++ {
		if !compareStatics(rends[i].Static, rends[i-1].Static) {
			staticsMatch = false
			break
		}
	}

	inserts := []any{}
	for _, r := range s.stream.Additions {
		inserts = append(inserts, []any{
			r.DomID,
			r.StreamAt,
			r.Limit,
		})
	}

	if staticsMatch {
		t.AddDynamic(&rend.Comprehension{
			Static:      rends[0].Static,
			Fingerprint: rends[0].Fingerprint,
			Dynamics:    copyDynamics(rends),
			Stream: []any{
				root.NextStreamID(),
				inserts,
				s.stream.Deletions,
				s.stream.Reset,
			},
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
