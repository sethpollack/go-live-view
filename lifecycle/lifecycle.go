package lifecycle

import (
	"fmt"

	"github.com/rs/xid"
	"github.com/sethpollack/go-live-view/html"
	lv "github.com/sethpollack/go-live-view/liveview"
	"github.com/sethpollack/go-live-view/params"
	"github.com/sethpollack/go-live-view/rend"
)

type Route interface {
	GetView() lv.LiveView
	GetParams() params.Params
}

type Router interface {
	GetRoute(string) (Route, error)
	Routable(Route, Route) bool
	GetLayout() func(...rend.Node) rend.Node
}

type lifecycle struct {
	router Router
	route  Route
	tree   *rend.Root
}

func NewLifecycle(r Router) *lifecycle {
	return &lifecycle{
		router: r,
	}
}

func (l *lifecycle) Join(s lv.Socket, p params.Params) (*rend.Root, error) {
	url := p.String("url", "redirect")

	route, err := l.router.GetRoute(url)
	if err != nil {
		return nil, err
	}

	l.route = route

	view := route.GetView()

	err = view.Mount(s, p)
	if err != nil {
		return nil, err
	}

	err = view.Params(s, nil)
	if err != nil {
		return nil, err
	}

	if s.Redirected() {
		return nil, nil
	}

	node, err := view.Render(nil)
	if err != nil {
		return nil, err
	}

	l.tree = rend.RenderTree(node)

	return l.tree, nil
}

func (l *lifecycle) Params(s lv.Socket, p params.Params) (*rend.Root, error) {
	url := p.String("url", "redirect")

	route, err := l.router.GetRoute(url)
	if err != nil {
		return nil, err
	}

	if !l.router.Routable(l.route, route) {
		err := s.Redirect(url)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("cant patch to %s, redirecting", url)

	}

	l.route = route

	view := route.GetView()

	if err := view.Params(s, p); err != nil {
		return nil, err
	}

	if s.Redirected() {
		return nil, nil
	}

	node, err := view.Render(nil)
	if err != nil {
		return nil, err
	}

	newTree := rend.RenderTree(node)

	diff := l.tree.Diff(newTree)

	l.tree = newTree

	return diff, nil
}

func (l *lifecycle) Event(s lv.Socket, p params.Params) (*rend.Root, error) {
	event := p.String("event")

	view := l.route.GetView()

	if err := view.Event(s, event, p); err != nil {
		return nil, err
	}

	if s.Redirected() {
		return nil, nil
	}

	node, err := view.Render(nil)
	if err != nil {
		return nil, err
	}

	newTree := rend.RenderTree(node)

	diff := l.tree.Diff(newTree)

	l.tree = newTree

	return diff, nil
}

func (l *lifecycle) StaticRender(sessionID string, url string) (string, error) {
	route, err := l.router.GetRoute(url)
	if err != nil {
		return "", err
	}

	view := route.GetView()

	err = view.Mount(nil, nil)
	if err != nil {
		return "", err
	}

	if err := view.Params(nil, nil); err != nil {
		return "", err
	}

	node, err := view.Render(nil)
	if err != nil {
		return "", err
	}

	return rend.RenderString(
		l.router.GetLayout()(
			html.Attrs(
				html.DataAttr("phx-main"),
				html.DataAttr("phx-session"), // TODO: do we need to implement this? Not sure exactly how this works.
				html.IdAttr(
					fmt.Sprintf("phx-%s", xid.New().String()),
				),
			),
			node,
		),
	), nil
}

func (l *lifecycle) DestroyCIDs(cids []int) error {
	if l.tree.Components != nil {
		for _, cid := range cids {
			if l.tree.Components[int64(cid)] != nil {
				return fmt.Errorf("component with cid %d found", cid)
			}
		}
	}

	return nil
}

func (l *lifecycle) Leave() error {
	return l.route.GetView().Unmount()
}

func (l *lifecycle) AllowUpload(s lv.Socket, p params.Params) (any, error) {
	ref := p.String("ref")

	view := l.route.GetView()

	u := view.Uploads()
	if u == nil {
		return nil, fmt.Errorf("uploads not found")
	}

	cfg := u.GetByRef(ref)
	if cfg == nil {
		return nil, fmt.Errorf("config not found")
	}

	cfg.OnAllowUploads(p)

	node, err := view.Render(nil)
	if err != nil {
		return nil, err
	}

	newTree := rend.RenderTree(node)

	diff := l.tree.Diff(newTree)

	l.tree = newTree

	if cfg == nil {
		return map[string]any{
			"diff": diff,
		}, nil
	}

	errors := cfg.PreflightErrors()
	if len(errors) > 0 {
		return map[string]any{
			"error": errors,
			"ref":   cfg.Ref,
		}, nil
	}

	return map[string]any{
		"config": map[string]any{
			"chunk_size":    cfg.ChunkSize,
			"max_entries":   cfg.MaxEntries,
			"max_file_size": cfg.MaxFileSize,
		},
		"diff":    diff,
		"entries": cfg.PreflightEntries(),
		"errors":  map[string]any{},
		"ref":     cfg.Ref,
	}, nil
}

func (l *lifecycle) Chunk(cRef, ref string, data []byte, close func() error) error {
	view := l.route.GetView()

	u := view.Uploads()
	if u == nil {
		return fmt.Errorf("uploads not found")
	}

	cfg := u.GetByRef(cRef)
	if cfg == nil {
		return fmt.Errorf("config not found")
	}

	return cfg.OnChunk(ref, data, close)
}

func (l *lifecycle) Progress(s lv.Socket, p params.Params) (*rend.Root, error) {
	ref := p.String("ref")
	eRef := p.String("entry_ref")
	progress := p.Float32("progress")

	view := l.route.GetView()

	u := view.Uploads()
	if u == nil {
		return nil, fmt.Errorf("uploads not found")
	}

	cfg := u.GetByRef(ref)
	if cfg == nil {
		return nil, fmt.Errorf("config not found")
	}

	err := cfg.OnProgress(eRef, progress)
	if err != nil {
		return nil, err
	}

	node, err := view.Render(nil)
	if err != nil {
		return nil, err
	}

	newTree := rend.RenderTree(node)

	diff := l.tree.Diff(newTree)

	l.tree = newTree

	return diff, nil
}
