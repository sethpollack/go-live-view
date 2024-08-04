package rend

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"go-live-view/internal/ref"
)

type Root struct {
	refCID    *ref.Ref
	streamRef *ref.Ref

	Components map[int64]*Rend `json:"c,omitempty"`
	Title      string          `json:"t,omitempty"`
	Rend       *Rend           `json:",inline"`
}

type Rend struct {
	refID *ref.Ref

	Dynamic     map[string]any `json:",inline"`
	Static      []string       `json:"s,omitempty"`
	Root        *bool          `json:"r,omitempty"`
	Fingerprint string         `json:"f,omitempty"`
}

type Comprehension struct {
	Static      []string `json:"s,omitempty"`
	Dynamics    [][]any  `json:"d,omitempty"`
	Fingerprint string   `json:"f,omitempty"`
	Stream      []any    `json:"stream,omitempty"`
}

func NewRoot() *Root {
	return &Root{
		refCID:    ref.New(0),
		streamRef: ref.New(-1),
		Rend:      &Rend{},
	}
}

func (r *Root) NextStreamID() int64 {
	return r.streamRef.NextRef()
}

func (rend *Rend) AddComponent(r *Root, c *Rend) {
	id := r.refCID.NextRef()

	rend.AddDynamic(id)

	if r.Components == nil {
		r.Components = map[int64]*Rend{}
	}

	c.Root = boolPtr(true)

	r.Components[id] = c
}

func (r *Rend) NextID() int64 {
	if r.refID == nil {
		r.refID = ref.New(-1)
	}
	return r.refID.NextRef()
}

func (r *Rend) AddStatic(s string) {
	r.Static = append(r.Static, s)
	r.Fingerprint = fingerPrint(r.Static)
}

func (r *Rend) AddDynamic(d any) {
	if r.Dynamic == nil {
		r.Dynamic = map[string]any{}
	}
	r.Dynamic[fmt.Sprintf("%d", r.NextID())] = d
}

func fingerPrint(s []string) string {
	h := sha256.New()
	for _, v := range s {
		h.Write([]byte(
			fmt.Sprintf("%s,", v), // single string vs multiple strings should have diff fingerprint
		))
	}
	return hex.EncodeToString(h.Sum(nil))
}

func boolPtr(b bool) *bool {
	return &b
}
