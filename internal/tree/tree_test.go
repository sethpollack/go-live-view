package tree

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindNode(t *testing.T) {
	tt := []struct {
		name         string
		paths        []string
		search       string
		expect       string
		expectParams map[string]any
		err          error
	}{
		{
			name:         "simple",
			paths:        []string{"/", "/foo", "/bar"},
			search:       "/",
			expect:       "/",
			expectParams: map[string]any{},
		},
		{
			name:         "partial match",
			paths:        []string{"/foo"},
			search:       "/foo/bar",
			expect:       "/foo",
			expectParams: map[string]any{},
		},
		{
			name:         "simple",
			paths:        []string{"/", "/foo", "/bar"},
			search:       "/bar",
			expect:       "/bar",
			expectParams: map[string]any{},
		},
		{
			name:         "simple with param",
			paths:        []string{"/:id"},
			search:       "/123",
			expect:       "/:id",
			expectParams: map[string]any{"id": "123"},
		},
		{
			name:         "nested params",
			paths:        []string{"/:id/:name"},
			search:       "/123/test",
			expect:       "/:id/:name",
			expectParams: map[string]any{"id": "123", "name": "test"},
		},
		{
			name:         "conflicting routes",
			paths:        []string{"/foo/123", "/foo/:id", "/foo/*"},
			search:       "/foo/123",
			expect:       "/foo/123",
			expectParams: map[string]any{},
		},
		{
			name:         "conflicting routes",
			paths:        []string{"/foo/123", "/foo/123/*", "/foo/:id/*"},
			search:       "/foo/123/bar",
			expect:       "/foo/123/*",
			expectParams: map[string]any{"*": "bar"},
		},
		{
			name:   "conflicting routes",
			paths:  []string{"/foo/123", "/foo/:id/*"},
			search: "/foo/123/bar",
			expect: "/foo/:id/*",
			expectParams: map[string]any{
				"id": "123",
				"*":  "bar",
			},
		},
		{
			name:   "duplicate route",
			paths:  []string{"/foo", "/foo"},
			search: "/foo",
			err:    errors.New("route /foo already exists"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tree := New[string]()
			for _, path := range tc.paths {
				err := tree.AddRoute(path, path)
				if err != nil {
					assert.Equal(t, tc.err, err)
					return
				}
			}

			node, params, err := tree.FindNode(tc.search)

			assert.Equal(t, tc.err, err)
			assert.Equal(t, tc.expect, node.GetRoute())
			assert.Equal(t, tc.expectParams, params)

		})
	}
}
