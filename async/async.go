package async

import (
	lv "github.com/sethpollack/go-live-view/liveview"
)

type State int

const (
	Loading State = iota
	Loaded
	Failed
)

type Async[T any] struct {
	value T
	state State
	err   error
}

func New[T any](s lv.Socket, fetch func() (T, error)) *Async[T] {
	a := &Async[T]{
		state: Loading,
	}

	if s == nil {
		return a
	}

	go func(s lv.Socket) {
		result, err := fetch()
		if err != nil {
			a.state = Failed
			a.err = err
		}
		a.value = result
		a.state = Loaded
		s.PushSelf("async:update", nil)
	}(s)

	return a
}

func (a *Async[T]) Value() T {
	return a.value
}

func (a *Async[T]) State() State {
	return a.state
}

func (a *Async[T]) Error() error {
	return a.err
}
