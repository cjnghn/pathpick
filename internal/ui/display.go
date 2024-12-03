package ui

import (
	"fmt"

	"github.com/cjnghn/pathpick/internal/tree"
)

type Display struct {
	root    *tree.Node
	current *tree.Node
	walker  *tree.Walker
}

func NewDisplay(walker *tree.Walker) *Display {
	return &Display{walker: walker}
}

func (d *Display) Start(rootPath string) error {
	var err error
	d.root, err = d.walker.Walk(rootPath)
	if err != nil {
		return err
	}
	d.current = d.root

	return d.eventLoop()
}

func (d *Display) render() {
	clearScreen()
	fmt.Println("PathPick - 파일 트리 탐색")
	fmt.Println("────────────────────────")
	fmt.Printf("현재 경로: %s\n", d.current.Path)
	fmt.Println("────────────────────────")
	fmt.Println("조작법:")
	fmt.Println("↑/↓: 이동  ←/→: 상위/하위  Space: 선택  y: 복사  q: 종료")
	fmt.Println("────────────────────────")

	d.renderNode(d.current, "", true)
}

func (d *Display) renderNode(node *tree.Node, prefix string, isLast bool) {
	marker := "├── "
	if isLast {
		marker = "└── "
	}

	selected := " "
	if node.Selected {
		selected = "*"
	}

	cursor := " "
	if node == d.current {
		cursor = ">"
	}

	fmt.Printf("%s%s%s%s %s\n", prefix, marker, cursor, selected, node.Name)

	childPrefix := prefix + "│   "
	if isLast {
		childPrefix = prefix + "    "
	}

	for i, child := range node.Children {
		d.renderNode(child, childPrefix, i == len(node.Children)-1)
	}
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}
