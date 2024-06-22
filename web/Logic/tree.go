package Logic

import (
	"strings"
)

type TreeNode struct {
	Name     string
	Path     string
	Children []*TreeNode
	IsEnd    bool
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
			isEnd := false
			if index == len(strs)-1 {
				isEnd = true
			}
			node := &TreeNode{
				Name:     name,
				Children: make([]*TreeNode, 0),
				Path:     t.Path + "/" + name,
				IsEnd:    isEnd,
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
	for index, name := range strs {
		if index == 0 {
			continue
		}
		children := t.Children
		var matchNode *TreeNode
		isMatch := false

		// 1. First, match the exact path
		for _, node := range children {
			if node.Name == name {
				matchNode = node
				isMatch = true
				break
			}
		}

		// 2. Second, match the path parameter
		if !isMatch {
			for _, node := range children {
				if strings.HasPrefix(node.Name, ":") {
					matchNode = node
					isMatch = true
					break
				}
			}
		}

		// 3. Third, match the single wildcard
		if !isMatch {
			for _, node := range children {
				if node.Name == "*" {
					matchNode = node
					isMatch = true
					break
				}
			}
		}

		// 4. Lastly, match the double wildcard
		if !isMatch {
			for _, node := range children {
				if node.Name == "**" {
					return node
				}
			}
		}

		if isMatch {
			t = matchNode
			if index == len(strs)-1 {
				return matchNode
			}
		} else {
			return nil
		}
	}
	return nil
}
