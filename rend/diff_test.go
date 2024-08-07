package rend

import (
	"flag"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	update = flag.Bool("update", false, "update .json files")
)

func TestDiff(t *testing.T) {
	tt := []struct {
		name string
		a    *Root
		b    *Root
	}{
		{
			name: "statics changed",
			a: &Root{
				Rend: &Rend{
					Static:      []string{"a", "b", "c"},
					Fingerprint: "123",
				},
			},
			b: &Root{
				Rend: &Rend{
					Static:      []string{"a", "b", "c", "d"},
					Fingerprint: "1234",
				},
			},
		},
		{
			name: "nested statics changed",
			a: &Root{
				Rend: &Rend{
					Dynamic: map[string]interface{}{
						"0": &Rend{
							Static:      []string{"a", "b", "c"},
							Fingerprint: "123",
						},
					},
				},
			},
			b: &Root{
				Rend: &Rend{
					Dynamic: map[string]interface{}{
						"0": &Rend{
							Static:      []string{"a", "b", "c", "d"},
							Fingerprint: "1234",
						},
					},
				},
			},
		},
		{
			name: "dynamics inserted",
			a: &Root{
				Rend: &Rend{},
			},
			b: &Root{
				Rend: &Rend{
					Dynamic: map[string]interface{}{
						"0": "a",
					},
				},
			},
		},
		{
			name: "dynamics changed",
			a: &Root{
				Rend: &Rend{
					Dynamic: map[string]interface{}{
						"0": &Rend{
							Dynamic: map[string]interface{}{
								"0": "a",
							},
						},
					},
				},
			},
			b: &Root{
				Rend: &Rend{
					Dynamic: map[string]interface{}{
						"0": &Rend{
							Dynamic: map[string]interface{}{
								"0": "b",
							},
						},
					},
				},
			},
		},
		{
			name: "dynamics added",
			a: &Root{
				Rend: &Rend{
					Dynamic: map[string]interface{}{
						"0": &Rend{
							Dynamic: map[string]interface{}{
								"0": "a",
							},
						},
					},
				},
			},
			b: &Root{
				Rend: &Rend{
					Dynamic: map[string]interface{}{
						"0": &Rend{
							Dynamic: map[string]interface{}{
								"0": "b",
							},
						},
						"1": &Rend{
							Dynamic: map[string]interface{}{
								"0": "b",
							},
						},
					},
				},
			},
		},
		{
			name: "dynamics removed",
			a: &Root{
				Rend: &Rend{
					Dynamic: map[string]interface{}{
						"0": &Rend{
							Dynamic: map[string]interface{}{
								"0": "a",
							},
						},
					},
				},
			},
			b: &Root{
				Rend: &Rend{},
			},
		},
		{
			name: "dynamic type changed",
			a: &Root{
				Rend: &Rend{
					Dynamic: map[string]interface{}{
						"0": "a",
					},
				},
			},
			b: &Root{
				Rend: &Rend{
					Dynamic: map[string]interface{}{
						"0": &Rend{
							Static:      []string{"", ""},
							Fingerprint: "123",
							Dynamic: map[string]interface{}{
								"0": "a",
							},
						},
					},
				},
			},
		},
		{
			name: "component added",
			a: &Root{
				Rend: &Rend{
					Fingerprint: "123",
				},
			},
			b: &Root{
				Rend: &Rend{
					Fingerprint: "123",
				},
				Components: map[int64]*Rend{
					1: {
						Static:      []string{"a", "b", "c"},
						Fingerprint: "123",
					},
				},
			},
		},
		{
			name: "component removed",
			a: &Root{
				Rend: &Rend{
					Fingerprint: "123",
				},
				Components: map[int64]*Rend{
					1: {
						Static:      []string{"a", "b", "c"},
						Fingerprint: "123",
					},
				},
			},
			b: &Root{
				Rend: &Rend{
					Fingerprint: "1234",
				},
			},
		},
		{
			name: "component updated",
			a: &Root{
				Rend: &Rend{
					Fingerprint: "1234",
				},
				Components: map[int64]*Rend{
					1: {
						Static:      []string{"a", "b", "c"},
						Fingerprint: "123",
					},
				},
			},
			b: &Root{
				Rend: &Rend{
					Fingerprint: "1234",
				},
				Components: map[int64]*Rend{
					1: {
						Static:      []string{"a", "b", "c", "d"},
						Fingerprint: "1234",
					},
				},
			},
		},
		{
			name: "comprehension statics changed",
			a: &Root{
				Rend: &Rend{
					Fingerprint: "123",
					Dynamic: map[string]interface{}{
						"0": &Comprehension{
							Static:      []string{"a", "b", "c"},
							Fingerprint: "123",
						},
					},
				},
			},
			b: &Root{
				Rend: &Rend{
					Fingerprint: "123",
					Dynamic: map[string]interface{}{
						"0": &Comprehension{
							Static:      []string{"a", "b", "c", "d"},
							Fingerprint: "1234",
						},
					},
				},
			},
		},
		{
			name: "comprehension dynamics changed",
			a: &Root{
				Rend: &Rend{
					Fingerprint: "123",
					Dynamic: map[string]interface{}{
						"0": &Comprehension{
							Fingerprint: "123",
							Dynamics: [][]any{
								{"a", "b", "c"},
							},
						},
					},
				},
			},
			b: &Root{
				Rend: &Rend{
					Fingerprint: "123",
					Dynamic: map[string]interface{}{
						"0": &Comprehension{
							Fingerprint: "123",
							Dynamics: [][]any{
								{"a", "b", "f"},
							},
						},
					},
				},
			},
		},
		{
			name: "comprehension dynamics added",
			a: &Root{
				Rend: &Rend{
					Fingerprint: "123",
					Dynamic: map[string]interface{}{
						"0": &Comprehension{
							Fingerprint: "123",
							Dynamics: [][]any{
								{"a", "b", "c"},
							},
						},
					},
				},
			},
			b: &Root{
				Rend: &Rend{
					Fingerprint: "123",
					Dynamic: map[string]interface{}{
						"0": &Comprehension{
							Fingerprint: "123",
							Dynamics: [][]any{
								{"a", "b", "c"},
								{"a", "b", "c"},
							},
						},
					},
				},
			},
		},
		{
			name: "comprehension component changed",
			a: &Root{
				Components: map[int64]*Rend{
					1: {
						Fingerprint: "123",
					},
				},
				Rend: &Rend{
					Fingerprint: "123",
					Dynamic: map[string]interface{}{
						"0": &Comprehension{
							Fingerprint: "123",
							Dynamics: [][]any{
								{int64(1)},
							},
						},
					},
				},
			},
			b: &Root{
				Components: map[int64]*Rend{
					1: {
						Fingerprint: "1234",
					},
				},
				Rend: &Rend{
					Fingerprint: "123",
					Dynamic: map[string]interface{}{
						"0": &Comprehension{
							Fingerprint: "123",
							Dynamics: [][]any{
								{int64(1)},
							},
						},
					},
				},
			},
		},
		{
			name: "comprehensions stream inserted",
			a: &Root{
				Rend: &Rend{
					Fingerprint: "123",
					Dynamic: map[string]interface{}{
						"0": &Comprehension{
							Fingerprint: "123",
							Stream:      []any{},
						},
					},
				},
			},
			b: &Root{
				Rend: &Rend{
					Fingerprint: "123",
					Dynamic: map[string]interface{}{
						"0": &Comprehension{
							Fingerprint: "123",
							Stream: []any{
								0, []any{"user-1", -1, nil}, []string{}, false,
							},
						},
					},
				},
			},
		},
		{
			name: "comprehensions stream deleted",
			a: &Root{
				Rend: &Rend{
					Fingerprint: "123",
					Dynamic: map[string]interface{}{
						"0": &Comprehension{
							Fingerprint: "123",
							Stream: []any{
								0, []any{"user-1", -1, nil}, []string{}, false,
							},
						},
					},
				},
			},
			b: &Root{
				Rend: &Rend{
					Fingerprint: "123",
					Dynamic: map[string]interface{}{
						"0": &Comprehension{
							Fingerprint: "123",
							Stream: []any{
								0, []any{}, []string{"user-1"}, false,
							},
						},
					},
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			diff := tc.a.Diff(tc.b)

			json := RenderJSONTree(diff)

			actual := actualValue(t, "testdata/"+stringify(tc.name)+".json", json, *update)

			assert.JSONEq(t, actual, json)
		})
	}
}

func actualValue(t *testing.T, path string, actual string, update bool) string {
	t.Helper()

	if update {
		err := os.WriteFile(path, []byte(actual), 0644)
		if err != nil {
			t.Fatalf("Error writing to file %s: %s", path, err)
		}
		return actual
	}

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Error opening file %s: %s", path, err)
	}

	return string(content)
}

func stringify(name string) string {
	return strings.ReplaceAll(name, " ", "-")
}
