package Logic

import (
	"strings"
)

type TreeNode struct {
	Name     string
	Path     string
	Children []*TreeNode
}

// put path: /user/get/:id
func (t *TreeNode) Put(path string) {
	strs := strings.Split(path, "/")
	for index, name := range strs {
		if index == 0 {
			continue
		}
		children := t.Children
		isMatch := false
		for _, node := range children {
			if node.Name == name {
				t = node
				isMatch = true
				break
			}
		}
		if !isMatch {
			node := &TreeNode{
				Name:     name,
				Children: make([]*TreeNode, 0),
			}
			children = append(children, node)
			t.Children = children
			t = node
		}
	}
}

// get path: /user/get/1
func (t *TreeNode) Get(path string) *TreeNode {
	strs := strings.Split(path, "/")
	fullPath := ""
	for index, name := range strs {
		if index == 0 {
			continue
		}
		children := t.Children
		isMatch := false
		for _, node := range children {
			if node.Name == name ||
				node.Name == "*" ||
				strings.Contains(node.Name, ":") {
				isMatch = true
				fullPath += "/" + node.Name
				node.Path = fullPath
				t = node
				if index == len(strs)-1 {
					return node
				}
				break
			}
		}
		if !isMatch {
			for _, node := range children {
				// /user/**
				if node.Name == "**" {
					return node
				}
			}
		}
	}
	return nil
}
