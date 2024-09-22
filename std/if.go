package std

import "github.com/sethpollack/go-live-view/rend"

func If(condition bool, then rend.Node) rend.Node {
	if condition {
		return DynamicNode(then)
	}
	return Noop()
}
