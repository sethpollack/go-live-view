package std

import (
	"github.com/sethpollack/go-live-view/rend"
)

func NotNil(cond any, cb func() rend.Node) rend.Node {
	if cond != nil {
		return DynamicNode(cb())
	}

	return DynamicNode(Noop())
}
