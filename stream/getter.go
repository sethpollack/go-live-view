package stream

import "fmt"

type StreamGetter struct {
	Stream *Stream
}

func New(name string, options ...StreamOption) *StreamGetter {
	return &StreamGetter{
		Stream: newStream(name, options...),
	}
}

func (s *StreamGetter) Get() *Stream {
	if s.Stream == nil {
		return nil
	}

	stream := s.Stream

	// Reset the stream
	s.Stream = &Stream{
		Name:      stream.Name,
		IDFunc:    stream.IDFunc,
		Additions: make([]Item, 0),
		Deletions: make([]string, 0),
		Reset:     stream.Reset,
		Limit:     stream.Limit,
		StreamAt:  stream.StreamAt,
	}

	return stream
}

func (s *StreamGetter) Add(items ...any) error {
	if s.Stream == nil {
		return fmt.Errorf("stream is nil")
	}

	s.Stream.add(items...)

	return nil
}

func (s *StreamGetter) Delete(ids ...string) error {
	if s.Stream == nil {
		return fmt.Errorf("stream is nil")
	}

	s.Stream.delete(ids...)

	return nil
}

func (s *StreamGetter) ResetStream() error {
	if s.Stream == nil {
		return fmt.Errorf("stream is nil")
	}

	s.Stream.resetStream()

	return nil
}
