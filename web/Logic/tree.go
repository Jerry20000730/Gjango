package Logic

import (
	"strings"
)

// TreeNode represents a node in a tree structure where each node can have multiple children.
// It is used to store and navigate hierarchical data, such as paths in a routing system.
type TreeNode struct {
	Name     string      // Name of the node, representing a segment in the path.
	Path     string      // Path is the full path from the root to this node.
	Children []*TreeNode // Children is a slice of pointers to child nodes.
	IsEnd    bool        // IsEnd indicates whether this node represents the end of a complete path.
}

// Put inserts a new path into the tree. If the path already exists, it ensures the path is marked as an endpoint.
// The path is split into segments, and a node is created for each segment if it doesn't already exist.
// Parameters:
// - path: The path to insert into the tree, starting with a slash (/).
func (t *TreeNode) Put(path string) {
	strs := strings.Split(path, "/") // Split the path into segments.
	for index, name := range strs {
		if index == 0 { // Skip the first segment if it's empty (leading slash).
			continue
		}
		children := t.Children
		isMatch := false
		for _, node := range children { // Look for an existing node with the same name.
			if node.Name == name {
				t = node // Move to the matching node.
				isMatch = true
				break
			}
		}
		if !isMatch { // If no matching node is found, create a new one.
			isEnd := false
			if index == len(strs)-1 { // Mark as an endpoint if it's the last segment.
				isEnd = true
			}
			node := &TreeNode{
				Name:     name,
				Children: make([]*TreeNode, 0),
				Path:     t.Path + "/" + name, // Construct the full path for the new node.
				IsEnd:    isEnd,
			}
			children = append(children, node) // Add the new node to the children of the current node.
			t.Children = children
			t = node // Move to the new node.
		}
	}
}

// Get retrieves a node from the tree that matches the given path.
// It supports matching exact paths, path parameters (prefixed with ":"), single wildcards (*),
// and double wildcards (**). The search priority is in the mentioned order.
// Parameters:
// - path: The path to search for in the tree, starting with a slash (/).
// Returns:
// - A pointer to the TreeNode that matches the path, or nil if no match is found.
func (t *TreeNode) Get(path string) *TreeNode {
	strs := strings.Split(path, "/") // Split the path into segments.
	for index, name := range strs {
		if index == 0 { // Skip the first segment if it's empty (leading slash).
			continue
		}
		children := t.Children
		var matchNode *TreeNode
		isMatch := false

		// 1. First, match the exact path.
		for _, node := range children {
			if node.Name == name {
				matchNode = node // Found a match with the exact name.
				isMatch = true
				break
			}
		}

		// 2. Second, match the path parameter (prefixed with ":").
		if !isMatch {
			for _, node := range children {
				if strings.HasPrefix(node.Name, ":") {
					matchNode = node // Found a match with a path parameter.
					isMatch = true
					break
				}
			}
		}

		// 3. Third, match the single wildcard (*).
		if !isMatch {
			for _, node := range children {
				if node.Name == "*" {
					matchNode = node // Found a match with a single wildcard.
					isMatch = true
					break
				}
			}
		}

		// 4. Lastly, match the double wildcard (**).
		if !isMatch {
			for _, node := range children {
				if node.Name == "**" {
					return node // Found a match with a double wildcard, return immediately.
				}
			}
		}

		if isMatch {
			t = matchNode             // Move to the matching node.
			if index == len(strs)-1 { // If it's the last segment, return the matching node.
				return matchNode
			}
		} else {
			return nil // No match found, return nil.
		}
	}
	return nil
}
