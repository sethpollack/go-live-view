package nested

import (
	lv "go-live-view/liveview"
	"go-live-view/params"
	"go-live-view/rend"
	"go-live-view/uploads"
)

var _ lv.LiveView = &view{}

type view struct {
	router *router
	route  *Route
}

func NewView(r *Route) *view {
	return &view{
		route:  r,
		router: r.router,
	}
}

func (v *view) Mount(s lv.Socket, p params.Params) error {
	return walk(v.route, func(route *Route) error {
		if !v.router.mounted[route] {
			v.router.mounted[route] = true
			return route.View.Mount(s, p)
		}
		return nil
	})
}

func (v *view) Unmount() error {
	for route := range v.router.mounted {
		err := route.View.Unmount()
		if err != nil {
			return err
		}
	}

	v.router.mounted = make(map[*Route]bool)

	return nil
}

func (v *view) Params(s lv.Socket, p params.Params) error {
	return walk(v.route, func(route *Route) error {
		if !v.router.mounted[route] {
			err := route.View.Mount(s, p)
			if err != nil {
				return err
			}
			v.router.mounted[route] = true
		}
		return route.View.Params(s, p)
	})
}

func (v *view) Event(s lv.Socket, e string, p params.Params) error {
	return walk(v.route, func(route *Route) error {
		return route.View.Event(s, e, p)
	})
}

func (v *view) Render(rend.Node) (node rend.Node, err error) {
	err = walk(v.route, func(route *Route) error {
		node, err = route.View.Render(node)
		return err
	})

	return node, err
}

func (v *view) Uploads() *uploads.Uploads {
	var u *uploads.Uploads

	walk(v.route, func(route *Route) error {
		if route.View != nil {
			uploads := route.View.Uploads()
			if uploads != nil {
				u = uploads
			}
		}
		return nil
	})

	return u
}

func walk(route *Route, f func(*Route) error) error {
	for route != nil {
		err := f(route)
		if err != nil {
			return err
		}
		route = route.parent
	}

	return nil
}
