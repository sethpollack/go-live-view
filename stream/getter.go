package stream

import "fmt"

type StreamGetter struct {
	opts   []StreamOption
	Stream *Stream
}

func New(name string, options ...StreamOption) *StreamGetter {
	return &StreamGetter{
		opts:   options,
		Stream: newStream(name, options...),
	}
}

func (s *StreamGetter) Get() *Stream {
	if s.Stream == nil {
		return nil
	}

	stream := s.Stream

	s.Stream = newStream(stream.Name, s.opts...)

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
