package tree

import (
	"os"
	"path/filepath"
)

type Walker struct {
	Pattern string
}

// 순회하면서 tree 완성
func (w *Walker) Walk(root string) (*Node, error) {
	info, err := os.Stat(root)
	if err != nil {
		return nil, err
	}

	node := &Node{
		Path:  root,
		Name:  filepath.Base(root),
		IsDir: info.IsDir(),
	}

	if !node.IsDir {
		if w.Pattern != "" {
			matched, err := filepath.Match(w.Pattern, node.Name)
			if err != nil || !matched {
				return nil, nil
			}
		}
		return node, nil
	}

	entries, err := os.ReadDir(root)
	if err != nil {
		return node, err
	}

	for _, entry := range entries {
		childPath := filepath.Join(root, entry.Name())
		child, err := w.Walk(childPath)
		if err != nil || child == nil {
			continue
		}

		child.Parent = node
		node.Children = append(node.Children, child)
	}

	// 마지막 자식 노드 표시
	if len(node.Children) > 0 {
		node.Children[len(node.Children)-1].IsLast = true
	}

	return node, nil
}
