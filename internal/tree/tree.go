package tree

import (
	"fmt"
	"strings"
)

type Node[T any] struct {
	segment  string
	children map[string]*Node[T]
	route    T
	depth    int
}

func New[T any]() *Node[T] {
	return &Node[T]{
		segment:  "",
		children: make(map[string]*Node[T]),
		depth:    1,
	}
}

func (n *Node[T]) GetRoute() T {
	return n.route
}

func (n *Node[T]) AddRoute(path string, route T) (*Node[T], error) {
	segments := strings.Split(path, "/")[1:]

	inserted := false

	current := n
	for _, segment := range segments {
		if child, ok := current.children[segment]; ok {
			current = child
			continue
		}

		newNode := &Node[T]{
			segment:  segment,
			children: make(map[string]*Node[T]),
			depth:    current.depth + 1,
		}

		current.children[segment] = newNode
		current = newNode
		inserted = true
	}

	if !inserted {
		return nil, fmt.Errorf("route %s already exists", path)
	}

	current.route = route

	return current, nil
}

func (n *Node[T]) FindNode(path string) (*Node[T], map[string]any, error) {
	segments := strings.Split(path, "/")[1:]

	currentFind := n
	currentParams := make(map[string]any)

	current := n
	for i, segment := range segments {
		if len(current.children) == 0 && i == len(segments)-1 && current.segment == segment {
			return current, currentParams, nil
		}

		currentPath := strings.Join(segments[i:], "/")

		if child, ok := current.children[segment]; ok {
			found, params, err := child.FindNode(currentPath)
			if err != nil {
				return nil, nil, err
			}

			if found.depth > currentFind.depth {
				currentParams = params
				currentFind = found
			}
		}

		for key, child := range current.children {
			if strings.HasPrefix(key, ":") {
				found, params, err := child.FindNode(currentPath)
				if err != nil {
					return nil, nil, err
				}

				if found.depth > currentFind.depth {
					params[key[1:]] = segment
					currentParams = params
					currentFind = found
				}
			}

			if key == "*" {
				found, params, err := child.FindNode(currentPath)
				if err != nil {
					return nil, nil, err
				}

				if found.depth > currentFind.depth {
					params["*"] = segment
					currentParams = params
					currentFind = found
				}
			}
			// if strings.HasPrefix(key, "~r") {
			// 	regex := key[2:]
			// 	r, err := regexp.Compile(regex)
			// 	if err != nil {
			// 		return nil, nil, err
			// 	}
			// 	if r.MatchString(segment) {
			// 		found, params, err := child.FindNode(currentPath)
			// 		if err != nil {
			// 			return nil, nil, err
			// 		}

			// 		if found.depth > currentFind.depth {
			// 			params[key] = segment
			// 			currentParams = params
			// 			currentFind = found
			// 		}
			// 	}
			// }
		}
	}

	return currentFind, currentParams, nil
}
