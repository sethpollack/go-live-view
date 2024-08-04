package ref

import (
	"fmt"
	"sync/atomic"
)

type Ref struct {
	ref *int64
}

func New(start int64) *Ref {
	return &Ref{
		ref: &start,
	}
}

func (r *Ref) NextRef() int64 {
	return atomic.AddInt64(r.ref, 1)
}

func (r *Ref) NextStringRef() string {
	return fmt.Sprintf("%d", r.NextRef())
}
