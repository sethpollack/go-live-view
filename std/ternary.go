package std

import "github.com/sethpollack/go-live-view/rend"

func TernaryString(cond bool, a, b string) *string {
	if cond {
		return &a
	}
	return &b
}

func TernaryNode(cond bool, a, b rend.Node) rend.Node {
	if cond {
		return DynamicNode(a)
	}
	return DynamicNode(b)
}

func TernaryNodeCB(cond bool, a, b func() rend.Node) rend.Node {
	if cond {
		return DynamicNode(a())
	}
	return DynamicNode(b())
}
