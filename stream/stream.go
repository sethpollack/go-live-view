package stream

type StreamOption func(*Stream)

func Reset() StreamOption {
	return func(s *Stream) {
		s.Reset = true
	}
}

func IDFunc(idFunc func(any) string) StreamOption {
	return func(s *Stream) {
		s.IDFunc = idFunc
	}
}

func Limit(limit int) StreamOption {
	return func(s *Stream) {
		s.Limit = &limit
	}
}

func StreamAt(streamAt int) StreamOption {
	return func(s *Stream) {
		s.StreamAt = &streamAt
	}
}

type Stream struct {
	Name string

	IDFunc func(any) string

	Additions []Item
	Deletions []string

	Reset    bool
	Limit    *int
	StreamAt *int
}

type Item struct {
	DomID    string
	Item     any
	StreamAt *int
	Limit    *int
}

func newStream(name string, options ...StreamOption) *Stream {
	s := &Stream{
		Name:      name,
		IDFunc:    func(any) string { panic("IDFunc not set") },
		Additions: make([]Item, 0),
		Deletions: make([]string, 0),
		Reset:     false,
	}

	for _, option := range options {
		option(s)
	}

	if s.StreamAt == nil {
		s.StreamAt = intPtr(-1)
	}

	return s
}

func (s *Stream) add(items ...any) {
	for _, item := range items {
		s.Additions = append(s.Additions, Item{
			DomID:    s.IDFunc(item),
			Item:     item,
			StreamAt: s.StreamAt,
			Limit:    s.Limit,
		})
	}
}

func (s *Stream) delete(ids ...string) {
	for _, id := range ids {
		s.Deletions = append(s.Deletions, id)
	}
}

func (s *Stream) resetStream() {
	s.Reset = true
}

func intPtr(i int) *int {
	return &i
}
