package html

import (
	"flag"
	"os"
	"strings"
	"testing"

	"github.com/sethpollack/go-live-view/rend"

	"github.com/stretchr/testify/assert"
)

var (
	update        = flag.Bool("update", false, "update .json files")
	dynamicText   = "hello"
	dynamicNumber = 123
)

func TestNode(t *testing.T) {
	tt := []struct {
		name  string
		tests []struct {
			name string
			node rend.Node
		}
	}{
		{
			name: "html attributes",
			tests: []struct {
				name string
				node rend.Node
			}{
				{
					name: "non attribute children are ignored for void elements",
					node: Void("div",
						Div(),
					),
				},
				{
					name: "simple void",
					node: Void("div",
						Attr("phx-click", "click"),
					),
				},
				{
					name: "attributes are ignored if wrapped in a dynamic",
					node: Div(
						dynamicNode(
							Attr("phx-click", "click"),
						),
					),
				},
				{
					name: "conditionals within attrs are not ignored",
					node: Div(
						Attrs(
							dynamicNode(
								Attr("phx-click", "click"),
							),
						),
					),
				},
				{
					name: "dynamic attributes",
					node: Div(
						Attr("attr1", &dynamicText),
						Attr("attr2", &dynamicNumber),
					),
				},
				{
					name: "dynamic attribute with empty string",
					node: Div(
						Attr("attr", strPtr("")),
					),
				},
				{
					name: "static attributes",
					node: Div(
						Attr("attr", "hello"),
						Attr("attr", 123),
					),
				},
				{
					name: "valueless attribute",
					node: Div(
						Attr("attr"),
					),
				},
				{
					name: "values helper",
					node: Div(
						Values("prefix-", map[string]any{
							"attr":  "hello",
							"other": 123,
						}),
					),
				},
			},
		},
		{
			name: "html elements",
			tests: []struct {
				name string
				node rend.Node
			}{
				{
					name: "nested elements",
					node: Html(
						Head(
							Meta(
								Title(),
							),
						),
						Body(
							Div(
								H1(),
							),
						),
					),
				},
				{
					name: "dynamic content",
					node: Div(
						dynamicNode(
							H1(),
						),
					),
				},
				{
					name: "nested dynamic content",
					node: Div(
						dynamicNode(
							dynamicNode(
								P(),
							),
						),
					),
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			for _, test := range tc.tests {
				t.Run(test.name, func(t *testing.T) {
					if test.node == nil {
						t.Skip("node is nil")
					}

					root := rend.RenderJSON(test.node)

					actual := actualValue(t, "testdata/"+stringify(tc.name+"-"+test.name)+".json", root, *update)

					assert.JSONEq(t, actual, root)
				})
			}
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

	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Error opening file %s: %s", path, err)
	}

	return string(b)
}

func stringify(name string) string {
	return strings.ReplaceAll(name, " ", "-")
}

type dNode struct {
	node rend.Node
}

func dynamicNode(root rend.Node) *dNode {
	return &dNode{root}
}

func (c *dNode) Render(diff bool, root *rend.Root, t *rend.Rend, b *strings.Builder) error {
	if diff {
		t.AddDynamic(rend.Render(root, c.node))
		t.AddStatic(b.String())
		b.Reset()

		return nil
	}

	return c.node.Render(diff, root, t, b)
}

func strPtr(s string) *string {
	return &s
}
