package tree

import (
	"os"
	"path/filepath"
	"strings"
)

type Walker struct {
	Pattern    string
	ShowHidden bool // 숨김 파일 표시 여부
}

func (w *Walker) shouldSkip(name string) bool {
	// 숨김 파일이고 ShowHidden이 false면 스킵
	if !w.ShowHidden && strings.HasPrefix(name, ".") {
		return true
	}

	// 패턴이 지정되어 있으면 패턴 매칭 검사
	if w.Pattern != "" {
		matched, err := filepath.Match(w.Pattern, name)
		if err != nil || !matched {
			return true
		}
	}

	return false
}

func (w *Walker) Walk(root string) (*Node, error) {
	info, err := os.Stat(root)
	if err != nil {
		return nil, err
	}

	// 루트 노드의 경우는 숨김 파일이어도 허용
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
		name := entry.Name()

		// 숨김 파일 및 패턴 검사
		if w.shouldSkip(name) {
			continue
		}

		childPath := filepath.Join(root, name)
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
