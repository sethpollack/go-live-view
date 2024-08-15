package channel

import (
	"net/http"
)

type Transport interface {
	Path() string
	Serve(func(Conn), http.ResponseWriter, *http.Request)
}
