package router

import (
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/sethpollack/go-live-view/internal/tree"
	lv "github.com/sethpollack/go-live-view/liveview"
	"github.com/sethpollack/go-live-view/params"
	"github.com/sethpollack/go-live-view/rend"
)

var _ lv.Route = (*route)(nil)
var _ lv.Router = (*router)(nil)

type routeOption func(*route)

type route struct {
	path   string
	view   lv.View
	params params.Params

	parent *route
	router *router
}

type router struct {
	root    *tree.Node[*route]
	mounted map[*route]bool
	layout  func(...rend.Node) rend.Node
}

type routeGroup struct {
	router  *router
	path    string
	options []routeOption
	parent  *route
}

func (r *route) GetView() lv.View {
	return newWrapper(r)
}

func (r *route) GetParams() params.Params {
	return r.params
}

func NewRouter(layout func(...rend.Node) rend.Node) *router {
	return &router{
		root:    tree.New[*route](),
		mounted: make(map[*route]bool),
		layout:  layout,
	}
}

func (r *router) GetLayout() func(...rend.Node) rend.Node {
	return r.layout
}

func (r *router) Group(path string, view lv.View, opts ...routeOption) *routeGroup {
	route := &route{
		path:   path,
		view:   view,
		router: r,
	}

	for _, opt := range opts {
		opt(route)
	}

	r.root.AddRoute(path, route)

	return &routeGroup{
		router:  r,
		path:    path,
		options: opts,
		parent:  route,
	}
}

func (r *router) Handle(path string, view lv.View, opts ...routeOption) *route {
	route := &route{
		path:   path,
		view:   view,
		router: r,
	}

	for _, opt := range opts {
		opt(route)
	}

	r.root.AddRoute(path, route)
	return route
}

func (r *router) GetRoute(path string) (lv.Route, error) {
	node, params, err := r.findNode(path)
	if err != nil {
		return nil, err
	}

	route := node.GetRoute()

	if route == nil {
		return nil, fmt.Errorf("no route found for path %s", path)
	}

	if route.params == nil {
		route.params = make(map[string]any)
	}

	for key, value := range params {
		route.params[key] = value
	}

	return route, nil
}

func (r *router) Routable(from lv.Route, to lv.Route) bool {
	return findParent(from.(*route)) == findParent(to.(*route))
}

func (r *router) findNode(path string) (*tree.Node[*route], map[string]any, error) {
	u, err := url.Parse(path)
	if err != nil {
		return nil, nil, err
	}

	node, params, err := r.root.FindNode(u.Path)
	if err != nil {
		return nil, nil, err
	}

	for key, value := range u.Query() {
		params[key] = value[0]
	}

	return node, params, nil
}

func (rg *routeGroup) Group(path string, view lv.View, opts ...routeOption) *routeGroup {
	fullPath := rg.combinePaths(rg.path, path)

	route := &route{
		path:   fullPath,
		view:   view,
		router: rg.router,
		parent: rg.parent,
	}

	for _, opt := range append(rg.options, opts...) {
		opt(route)
	}

	// Add the group's route to the router
	rg.router.root.AddRoute(fullPath, route)

	return &routeGroup{
		router:  rg.router,
		path:    fullPath,
		options: append(rg.options, opts...),
		parent:  route,
	}
}

func (rg *routeGroup) Handle(path string, view lv.View, opts ...routeOption) *route {
	fullPath := rg.combinePaths(rg.path, path)

	route := &route{
		path:   fullPath,
		view:   view,
		router: rg.router,
		parent: rg.parent,
	}

	for _, opt := range append(rg.options, opts...) {
		opt(route)
	}

	rg.router.root.AddRoute(fullPath, route)
	return route
}

func (rg *routeGroup) combinePaths(base, new string) string {
	new = strings.TrimPrefix(new, "/")
	return path.Join(base, new)
}

func WithParams(params params.Params) routeOption {
	return func(r *route) {
		combined := make(map[string]any)
		for k, v := range r.params {
			combined[k] = v
		}
		for k, v := range params {
			combined[k] = v
		}
		r.params = combined
	}
}

func findParent(route *route) *route {
	if route == nil {
		return nil
	}

	for route != nil && route.parent != nil {
		route = route.parent
	}

	return route
}
