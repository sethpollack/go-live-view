package liveview

import "net/http"

type sessionGetter interface {
	Get(*http.Request) map[string]any
}
