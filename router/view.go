package router

import (
	lv "github.com/sethpollack/go-live-view/liveview"
	"github.com/sethpollack/go-live-view/params"
	"github.com/sethpollack/go-live-view/rend"
	"github.com/sethpollack/go-live-view/uploads"
)

var _ interface {
	lv.View
	lv.Mounter
	lv.Unmounter
	lv.Patcher
	lv.EventHandler
	lv.Uploader
} = &wrapper{}

type wrapper struct {
	router *router
	route  *route
}

func newWrapper(r *route) *wrapper {
	return &wrapper{
		route:  r,
		router: r.router,
	}
}

func (v *wrapper) Mount(s lv.Socket, p params.Params) error {
	return walk(v.route, func(route *route) error {
		if !v.router.mounted[route] {
			v.router.mounted[route] = true
			return lv.TryMount(route.view, s, p)
		}
		return nil
	})
}

func (v *wrapper) Unmount() error {
	for route := range v.router.mounted {
		err := lv.TryUnmount(route.view)
		if err != nil {
			return err
		}
	}

	v.router.mounted = make(map[*route]bool)

	return nil
}

func (v *wrapper) Params(s lv.Socket, p params.Params) error {
	// unmount everything except the current route
	if v.route.parent == nil {
		for current := range v.router.mounted {
			if current != v.route {
				lv.TryUnmount(current.view)
				delete(v.router.mounted, current)
			}
		}
	}

	// mount everything that is in the current tree and track it.
	shouldMount := make(map[*route]bool)
	err := walk(v.route, func(route *route) error {
		if !v.router.mounted[route] {
			err := lv.TryMount(route.view, s, p)
			if err != nil {
				return err
			}
			v.router.mounted[route] = true
		}
		shouldMount[route] = true
		return lv.TryParams(route.view, s, p)
	})

	// unmount anything that is not in the current tree.
	for mount := range v.router.mounted {
		if !shouldMount[mount] {
			lv.TryUnmount(mount.view)
			delete(v.router.mounted, mount)
		}
	}

	return err
}

func (v *wrapper) Event(s lv.Socket, e string, p params.Params) error {
	return walk(v.route, func(route *route) error {
		return lv.TryEvent(route.view, s, e, p)
	})
}

func (v *wrapper) Render(rend.Node) (node rend.Node, err error) {
	err = walk(v.route, func(route *route) error {
		node, err = route.view.Render(node)
		return err
	})

	return node, err
}

func (v *wrapper) Uploads() *uploads.Uploads {
	var u *uploads.Uploads

	walk(v.route, func(route *route) error {
		if route.view != nil {
			uploads := lv.TryUploads(route.view)
			if uploads != nil {
				u = uploads
			}
		}
		return nil
	})

	return u
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
