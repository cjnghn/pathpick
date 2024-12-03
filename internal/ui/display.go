package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cjnghn/pathpick/internal/tree"
	"golang.org/x/term"
)

const logo = `
██████╗  █████╗ ████████╗██╗  ██╗██████╗ ██╗ ██████╗██╗  ██╗
██╔══██╗██╔══██╗╚══██╔══╝██║  ██║██╔══██╗██║██╔════╝██║ ██╔╝
██████╔╝███████║   ██║   ███████║██████╔╝██║██║     █████╔╝
██╔═══╝ ██╔══██║   ██║   ██╔══██║██╔═══╝ ██║██║     ██╔═██╗
██║     ██║  ██║   ██║   ██║  ██║██║     ██║╚██████╗██║  ██╗
╚═╝     ╚═╝  ╚═╝   ╚═╝   ╚═╝  ╚═╝╚═╝     ╚═╝ ╚═════╝╚═╝  ╚═╝
`

const divider = `┌────────────────────────────────────────────────────────────┐`
const dividerMid = `├────────────────────────────────────────────────────────────┤`
const dividerEnd = `└────────────────────────────────────────────────────────────┘`

type Display struct {
	root    *tree.Node
	current *tree.Node
	walker  *tree.Walker
	offset  int
	height  int
}

func NewDisplay(walker *tree.Walker) *Display {
	height, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		height = 24 // 기본값
	}
	return &Display{
		walker: walker,
		height: height,
	}
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

func (d *Display) renderHeader() int {
	lines := 0

	// Logo
	fmt.Println(logo)
	lines += strings.Count(logo, "\n") + 1

	fmt.Println(divider)
	lines++

	selected := d.root.GetSelected()
	if len(selected) > 0 {
		fmt.Println("│ Selected Files:")
		lines++
		for i, node := range selected {
			relPath, err := filepath.Rel(d.root.Path, node.Path)
			if err != nil {
				relPath = node.Path
			}
			if i < 3 {
				fmt.Printf("│  • %s\n", relPath)
				lines++
			} else if i == 3 {
				fmt.Printf("│  • ... and %d more file(s)\n", len(selected)-3)
				lines++
				break
			}
		}
		fmt.Println(dividerMid)
		lines++
	}

	fmt.Printf("│ Current: %s\n", d.current.Path)
	fmt.Printf("│ Total selected: %d file(s)\n", len(selected))
	fmt.Println(dividerMid)
	lines += 3

	fmt.Println("│ Navigation:")
	fmt.Println("│  ↑/↓ : Move cursor    ←/→ : Parent/Child")
	fmt.Println("│ Space: Toggle select    y  : Copy to clipboard")
	fmt.Println("│    q : Quit            h  : Toggle hidden files")
	fmt.Println(dividerEnd)
	lines += 5

	return lines
}

func (d *Display) renderNodeToLines(node *tree.Node, prefix string, isLast bool, lines *[]string) {
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

	name := node.Name
	if node.IsDir {
		name = "[" + name + "]"
	}

	*lines = append(*lines, fmt.Sprintf("%s%s%s%s %s", prefix, marker, cursor, selected, name))

	childPrefix := prefix + "│   "
	if isLast {
		childPrefix = prefix + "    "
	}

	for i, child := range node.Children {
		d.renderNodeToLines(child, childPrefix, i == len(node.Children)-1, lines)
	}
}

func (d *Display) render() {
	clearScreen()

	// 헤더 렌더링
	headerLines := d.renderHeader()

	// 트리 렌더링
	lines := []string{}
	d.renderNodeToLines(d.root, "", true, &lines)

	// 현재 노드의 위치 찾기
	currentLine := d.findCurrentNodeLine(lines)

	// 화면에 맞게 오프셋 조정
	availableLines := d.height - headerLines
	if currentLine < d.offset {
		d.offset = currentLine
	} else if currentLine >= d.offset+availableLines {
		d.offset = currentLine - availableLines + 1
	}

	// 페이지에 맞는 라인만 출력
	endLine := min(d.offset+availableLines, len(lines))
	for _, line := range lines[d.offset:endLine] {
		fmt.Println(line)
	}

	// 스크롤 정보 표시
	if len(lines) > availableLines {
		fmt.Println(divider)
		fmt.Printf("│ Page: %d/%d\n", d.offset+availableLines, len(lines))
		fmt.Println(dividerEnd)
	}
}

func (d *Display) findCurrentNodeLine(lines []string) int {
	for i, line := range lines {
		if strings.Contains(line, "> ") {
			return i
		}
	}
	return 0
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}
