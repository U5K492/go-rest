package router

import (
	"strings"
)

type TreeNode struct {
	children []*TreeNode
	handler  func(ctx *Context)
	param    string
	parent   *TreeNode
}

func Constructor() *TreeNode {
	return &TreeNode{
		param:    "",
		children: []*TreeNode{},
	}
}

func (t *TreeNode) Insert(pathName string, handler func(ctx *Context)) {
	node := t

	params := strings.Split(pathName, "/")

	for _, param := range params {
		child := node.findChild(param)

		if child == nil {
			child = &TreeNode{
				param:    param,
				children: []*TreeNode{},
				parent:   node,
			}

			node.children = append(node.children, child)
		}

		node = child
	}
	node.handler = handler
}

func (t *TreeNode) findChild(param string) *TreeNode {
	for _, child := range t.children {
		if child.param == param {
			return child
		}
	}
	return nil
}

func (t *TreeNode) walk(pathName string) *TreeNode {
	params := strings.Split(pathName, "/")

	return dfs(t, params)
}

func dfs(node *TreeNode, params []string) *TreeNode {
	currentParam := params[0]
	isLastParam := len(params) == 1

	for _, child := range node.children {

		if isLastParam {
			if strings.HasPrefix(child.param, ":") {
				return child
			}

			if child.param == currentParam {
				return child
			}

			continue
		}

		if !strings.HasPrefix(child.param, ":") && child.param != currentParam {
			continue
		}

		result := dfs(child, params[1:])

		if result != nil {
			return result
		}
	}

	return nil
}

func (t *TreeNode) ParseParams(pathName string) map[string]string {
	node := t

	params := strings.Split(pathName, "/")

	paramDict := make(map[string]string)

	for i := len(params) - 1; i >= 0; i-- {
		if strings.HasPrefix(node.param, ":") {
			paramDict[node.param] = params[i]
		}
		node = node.parent
	}

	return paramDict
}
