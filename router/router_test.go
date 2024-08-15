package router

import (
	"testing"

	"github.com/sethpollack/go-live-view/html"
	lv "github.com/sethpollack/go-live-view/liveview"
	"github.com/sethpollack/go-live-view/params"
	"github.com/sethpollack/go-live-view/rend"
	"github.com/sethpollack/go-live-view/std"

	"github.com/stretchr/testify/assert"
)

func testLayout(name string, n rend.Node) rend.Node {
	return std.Noop()
}

type testLive struct {
	name string
	lv.Base
}

func (t *testLive) Render(n rend.Node) (rend.Node, error) {
	return html.Div(
		std.Text(t.name),
		n,
	), nil
}

type routes struct {
	path     string
	lv       lv.LiveView
	children []routes
	opts     []routeOption
}

type handler interface {
	Handle(path string, view lv.LiveView, opts ...routeOption) *route
	Group(path string, view lv.LiveView, opts ...routeOption) *routeGroup
}

func TestRouter(t *testing.T) {
	tt := []struct {
		name           string
		path           string
		routes         []routes
		expected       string
		expectedParams params.Params
	}{
		{
			name: "simple routes",
			path: "/",
			routes: []routes{
				{path: "/", lv: &testLive{
					name: "test",
				}},
			},
			expected:       "<div>test</div>",
			expectedParams: params.Params{},
		},
		{
			name: "simple routes with route params",
			path: "/123",
			routes: []routes{
				{path: "/:id", lv: &testLive{
					name: "test",
				}},
			},
			expected: "<div>test</div>",
			expectedParams: params.Params{
				"id": "123",
			},
		},
		{
			name: "simple routes with wildcard",
			path: "/test/123",
			routes: []routes{
				{path: "/test/*", lv: &testLive{
					name: "test",
				}},
			},
			expected: "<div>test</div>",
			expectedParams: params.Params{
				"*": "123",
			},
		},
		{
			name: "simple routes with query params",
			path: "/test?a=1&b=2",
			routes: []routes{
				{path: "/test", lv: &testLive{
					name: "test",
				}},
			},
			expected: "<div>test</div>",
			expectedParams: params.Params{
				"a": "1",
				"b": "2",
			},
		},
		{
			name: "simple routes with extra params",
			path: "/",
			routes: []routes{
				{path: "/", lv: &testLive{
					name: "test",
				},
					opts: []routeOption{
						WithParams(map[string]interface{}{
							"extra": "extra",
						}),
					},
				},
			},
			expected: "<div>test</div>",
			expectedParams: map[string]any{
				"extra": "extra",
			},
		},
		{
			name: "nested routes with route params",
			path: "/test/123",
			routes: []routes{
				{path: "/test", lv: &testLive{
					name: "test",
				},
					children: []routes{
						{path: "/:id", lv: &testLive{
							name: "child",
						}},
					},
				},
			},
			expected: "<div>test<div>child</div></div>",
			expectedParams: params.Params{
				"id": "123",
			},
		},
		{
			name: "nested routes with extra params",
			path: "/test/child/deep",
			routes: []routes{
				{
					path: "/test",
					lv: &testLive{
						name: "test",
					},
					opts: []routeOption{
						WithParams(map[string]interface{}{
							"extra1": "extra 1",
						}),
					},
					children: []routes{
						{
							path: "/child",
							lv: &testLive{
								name: "child",
							},
							opts: []routeOption{
								WithParams(map[string]interface{}{
									"extra2": "extra 2",
								}),
							},
							children: []routes{
								{
									path: "/deep",
									lv: &testLive{
										name: "deep",
									},
									opts: []routeOption{
										WithParams(map[string]interface{}{
											"extra3": "extra 3",
										}),
									},
								},
							},
						},
					},
				},
			},
			expected: "<div>test<div>child<div>deep</div></div></div>",
			expectedParams: map[string]any{
				"extra1": "extra 1",
				"extra2": "extra 2",
				"extra3": "extra 3",
			},
		},
		{
			name: "nested routes",
			path: "/test/child",
			routes: []routes{
				{path: "/test", lv: &testLive{
					name: "test",
				},
					children: []routes{
						{path: "/child", lv: &testLive{
							name: "child",
						}},
					},
				},
			},
			expected:       "<div>test<div>child</div></div>",
			expectedParams: params.Params{},
		},
		{
			name: "deep nested routes",
			path: "/test/child/deep",
			routes: []routes{
				{path: "/test", lv: &testLive{
					name: "test",
				},
					children: []routes{
						{path: "/child", lv: &testLive{
							name: "child",
						},
							children: []routes{
								{path: "/deep", lv: &testLive{
									name: "deep",
								}},
							},
						},
					},
				},
			},
			expected:       "<div>test<div>child<div>deep</div></div></div>",
			expectedParams: params.Params{},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			rt := NewRouter(
				testLayout,
			)
			createRoutes(rt, tc.routes)
			route, err := rt.GetRoute(tc.path)
			if err != nil {
				t.Fatalf("error getting route: %v", err)
			}

			assert.Equal(t, tc.expectedParams, route.GetParams())

			node, err := route.GetView().Render(nil)
			if err != nil {
				t.Fatalf("error rendering route: %v", err)
			}

			result := rend.RenderString(node)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func createRoutes(rt handler, routes []routes) {
	for _, route := range routes {
		if len(route.children) > 0 {
			createRoutes(rt.Group(route.path, route.lv, route.opts...), route.children)
		} else {
			rt.Handle(route.path, route.lv, route.opts...)
		}
	}
}
