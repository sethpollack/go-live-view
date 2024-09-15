package handler

import "net/http"

type sessionGetter interface {
	Get(*http.Request) map[string]any
}

type defaultSessionGetter struct{}

func (d *defaultSessionGetter) Get(r *http.Request) map[string]any {
	return make(map[string]any)
}
