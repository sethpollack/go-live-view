package router

import (
	"fmt"
	"go-live-view/internal/tree"
	"go-live-view/lifecycle"
	lv "go-live-view/liveview"
	"go-live-view/params"
	"go-live-view/rend"
	"net/url"
	"path"
	"strings"
)

var _ lifecycle.Route = (*route)(nil)
var _ lifecycle.Router = (*router)(nil)

type routeOption func(*route)

type route struct {
	Path   string
	View   lv.LiveView
	Params params.Params
	Layout func(string, rend.Node) rend.Node

	parent *route
	router *router
}

type router struct {
	root    *tree.Node[*route]
	mounted map[*route]bool
	options []routeOption
}

type routeGroup struct {
	router  *router
	path    string
	options []routeOption
	parent  *route
}

func (r *route) GetView() lv.LiveView {
	return NewWrapper(r)
}

func (r *route) GetParams() params.Params {
	return r.Params
}

func (r *route) GetLayout() func(string, rend.Node) rend.Node {
	return findParent(r).Layout
}

func NewRouter(opts ...routeOption) *router {
	return &router{
		root:    tree.New[*route](),
		mounted: make(map[*route]bool),
		options: opts,
	}
}

func (r *router) Group(path string, view lv.LiveView, opts ...routeOption) *routeGroup {
	route := &route{
		Path:   path,
		View:   view,
		router: r,
	}

	for _, opt := range append(r.options, opts...) {
		opt(route)
	}

	r.root.AddRoute(path, route)

	return &routeGroup{
		router:  r,
		path:    path,
		options: append(r.options, opts...),
		parent:  route,
	}
}

func (r *router) Handle(path string, view lv.LiveView, opts ...routeOption) *route {
	route := &route{
		Path:   path,
		View:   view,
		router: r,
	}

	for _, opt := range append(r.options, opts...) {
		opt(route)
	}

	r.root.AddRoute(path, route)
	return route
}

func (r *router) GetRoute(path string) (lifecycle.Route, error) {
	node, params, err := r.findNode(path)
	if err != nil {
		return nil, err
	}

	route := node.GetRoute()

	if route == nil {
		return nil, fmt.Errorf("no route found for path %s", path)
	}

	if route.Params == nil {
		route.Params = make(map[string]any)
	}

	for key, value := range params {
		route.Params[key] = value
	}

	return route, nil
}

func (r *router) Routable(from lifecycle.Route, to lifecycle.Route) bool {
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

func (rg *routeGroup) Group(path string, view lv.LiveView, opts ...routeOption) *routeGroup {
	fullPath := rg.combinePaths(rg.path, path)

	route := &route{
		Path:   fullPath,
		View:   view,
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

func (rg *routeGroup) Handle(path string, view lv.LiveView, opts ...routeOption) *route {
	fullPath := rg.combinePaths(rg.path, path)

	route := &route{
		Path:   fullPath,
		View:   view,
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

func WithLayout(layout func(string, rend.Node) rend.Node) routeOption {
	return func(r *route) {
		r.Layout = layout
	}
}

func WithParams(params params.Params) routeOption {
	return func(r *route) {
		combined := make(map[string]any)
		for k, v := range r.Params {
			combined[k] = v
		}
		for k, v := range params {
			combined[k] = v
		}
		r.Params = combined
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
