package std

import (
	"flag"
	"os"
	"strings"
	"testing"

	"github.com/sethpollack/go-live-view/html"
	"github.com/sethpollack/go-live-view/rend"

	"github.com/stretchr/testify/assert"
)

var (
	update        = flag.Bool("update", false, "update .json files")
	dynamicText   = "dynamic text"
	dynamicNumber = 456
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
			name: "components",
			tests: []struct {
				name string
				node rend.Node
			}{
				{
					name: "simple component",
					node: Component(
						Raw(`<div>Hello World</div>`),
					),
				},
				{
					name: "nested component",
					node: Component(
						Component(
							Raw(`<div>Hello World</div>`),
						),
					),
				},
			},
		},
		{
			name: "comprehensions",
			tests: []struct {
				name string
				node rend.Node
			}{
				{
					name: "simple comprehension",
					node: Range[string](
						[]string{"a", "b", "c"},
						func(s string) rend.Node {
							return Raw(`<div>Hello World</div>`)
						},
					),
				},
				{
					name: "changing statics comprehension",
					node: html.Div(
						Range[string](
							[]string{"a", "b", "c"},
							func(s string) rend.Node {
								return Raw(`<div>Hello World ` + s + ` </div>`)
							},
						),
					),
				},
				{
					name: "nested comprehension",
					node: Range[string](
						[]string{"a", "b", "c"},
						func(s string) rend.Node {
							return Range[string](
								[]string{"a", "b", "c"},
								func(s string) rend.Node {
									return Raw(`<div>Hello World</div>`)
								},
							)
						},
					),
				},
				{
					name: "component in comprehension",
					node: Range[string](
						[]string{"a", "b", "c"},
						func(s string) rend.Node {
							return Component(
								Raw(`<div>Hello World</div>`),
							)
						},
					),
				},
			},
		},
		{
			name: "group",
			tests: []struct {
				name string
				node rend.Node
			}{
				{
					name: "simple",
					node: html.Div(
						Group(
							Raw(`<div>Hello World 1</div>`),
							Raw(`<div>Hello World 2</div>`),
							Raw(`<div>Hello World 3</div>`),
						),
					),
				},
			},
		},
		{
			name: "noop",
			tests: []struct {
				name string
				node rend.Node
			}{
				{
					name: "simple",
					node: Noop(),
				},
			},
		},
		{
			name: "goembed",
			tests: []struct {
				name string
				node rend.Node
			}{
				{
					name: "simple",
					node: GoEmbed(
						func() rend.Node {
							return Raw(`<div>Hello World</div>`)
						},
					),
				},
			},
		},
		{
			name: "raw",
			tests: []struct {
				name string
				node rend.Node
			}{
				{
					name: "simple",
					node: Raw(`<div>Hello World</div>`),
				},
			},
		},
		{
			name: "notnil",
			tests: []struct {
				name string
				node rend.Node
			}{
				{
					name: "simple",
					node: NotNil(true, func() rend.Node {
						return Raw(`<div>Hello World</div>`)
					}),
				},
			},
		},
		{
			name: "ternary",
			tests: []struct {
				name string
				node rend.Node
			}{
				{
					name: "simple",
					node: TernaryNode(true,
						Raw(`<div>true</div>`),
						Raw(`<div>false</div>`),
					),
				},
				{
					name: "string",
					node: Raw(
						*TernaryString(true, "a", "b"),
					),
				},
				{
					name: "callback",
					node: html.Div(
						TernaryNodeCB(true,
							func() rend.Node {
								return Raw(`<div>true</div>`)
							},
							func() rend.Node {
								return Raw(`<div>false</div>`)
							},
						),
					),
				},
			},
		},
		{
			name: "if true",
			tests: []struct {
				name string
				node rend.Node
			}{
				{
					name: "simple",
					node: If(true, Raw(`<div>true</div>`)),
				},
			},
		},
		{
			name: "if false",
			tests: []struct {
				name string
				node rend.Node
			}{
				{
					name: "simple",
					node: If(false, Raw(`<div>true</div>`)),
				},
			},
		},
		{
			name: "text",
			tests: []struct {
				name string
				node rend.Node
			}{
				{
					name: "simple",
					node: Text("a"),
				},
				{
					name: "dynamic",
					node: Text(&dynamicText),
				},
				{
					name: "dynamic with format",
					node: Textf(
						"Hello %s %d",
						&dynamicText,
						&dynamicNumber,
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
					json := rend.RenderJSON(test.node)
					actual := actualValue(t, "testdata/"+stringify(tc.name+"-"+test.name)+".json", json, *update)

					assert.JSONEq(t, actual, json)
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

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Error opening file %s: %s", path, err)
	}

	return string(content)
}

func stringify(name string) string {
	return strings.ReplaceAll(name, " ", "-")
}
