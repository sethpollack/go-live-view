package router

import (
	lv "github.com/sethpollack/go-live-view/liveview"
	"github.com/sethpollack/go-live-view/params"
	"github.com/sethpollack/go-live-view/rend"
	"github.com/sethpollack/go-live-view/uploads"
)

var _ lv.LiveView = &wrapper{}

type wrapper struct {
	router *router
	route  *route
}

func NewWrapper(r *route) *wrapper {
	return &wrapper{
		route:  r,
		router: r.router,
	}
}

func (v *wrapper) Mount(s lv.Socket, p params.Params) error {
	return walk(v.route, func(route *route) error {
		if !v.router.mounted[route] {
			v.router.mounted[route] = true
			return route.View.Mount(s, p)
		}
		return nil
	})
}

func (v *wrapper) Unmount() error {
	for route := range v.router.mounted {
		err := route.View.Unmount()
		if err != nil {
			return err
		}
	}

	v.router.mounted = make(map[*route]bool)

	return nil
}

func (v *wrapper) Params(s lv.Socket, p params.Params) error {
	v.unmountSiblings(v.route)

	return walk(v.route, func(route *route) error {
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

func (v *wrapper) Event(s lv.Socket, e string, p params.Params) error {
	return walk(v.route, func(route *route) error {
		return route.View.Event(s, e, p)
	})
}

func (v *wrapper) Render(rend.Node) (node rend.Node, err error) {
	err = walk(v.route, func(route *route) error {
		node, err = route.View.Render(node)
		return err
	})

	return node, err
}

func (v *wrapper) Uploads() *uploads.Uploads {
	var u *uploads.Uploads

	walk(v.route, func(route *route) error {
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

func (v *wrapper) unmountSiblings(route *route) {
	if route.parent == nil {
		return
	}

	for sibling := range v.router.mounted {
		if sibling.parent == route.parent && sibling != route {
			sibling.View.Unmount()
			delete(v.router.mounted, sibling)
		}
	}
}

func walk(route *route, f func(*route) error) error {
	for route != nil {
		err := f(route)
		if err != nil {
			return err
		}
		route = route.parent
	}

	return nil
}
