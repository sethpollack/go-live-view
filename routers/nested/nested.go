package nested

import (
	"fmt"
	"go-live-view/internal/tree"
	"go-live-view/lifecycle"
	lv "go-live-view/liveview"
	"go-live-view/params"
	"go-live-view/rend"
	"net/url"
)

var _ lifecycle.Route = (*Route)(nil)

type Route struct {
	Path   string
	View   lv.LiveView
	Params params.Params
	Layout func(string, rend.Node) rend.Node

	parent *Route
	router *router
}

type Routes struct {
	Route    *Route
	Children []*Routes
}

func (r *Route) GetView() lv.LiveView {
	return NewView(r)
}

func (r *Route) GetParams() params.Params {
	return r.Params
}

func (r *Route) GetLayout() func(string, rend.Node) rend.Node {
	return findParent(r).Layout
}

type router struct {
	mounted map[*Route]bool
	root    *tree.Node[*Route]
}

func NewRouter() *router {
	return &router{
		mounted: make(map[*Route]bool),
		root:    tree.New[*Route](),
	}
}

func (r *router) HandleLive(routes *Routes) error {
	return r.addRoutes(routes.Route.Path, routes)
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
	return findParent(from.(*Route)) == findParent(to.(*Route))
}

func (r *router) findNode(path string) (*tree.Node[*Route], map[string]any, error) {
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

func (r *router) addRoutes(path string, routes *Routes) error {
	if routes.Route.Path == "" {
		return fmt.Errorf("path is required")
	}

	if routes.Route.View == nil {
		return fmt.Errorf("view is required")
	}

	routes.Route.router = r

	err := r.root.AddRoute(path, routes.Route)
	if err != nil {
		return err
	}

	for _, child := range routes.Children {
		if child.Route.Layout != nil {
			return fmt.Errorf("child routes do not support layout")
		}
		child.Route.parent = routes.Route
		var childPath string
		if path == "/" {
			childPath = child.Route.Path
		} else {
			childPath = path + child.Route.Path
		}

		err := r.addRoutes(childPath, child)
		if err != nil {
			return err
		}
	}

	return nil
}

func findParent(route *Route) *Route {
	if route == nil {
		return nil
	}

	for route != nil && route.parent != nil {
		route = route.parent
	}

	return route
}
