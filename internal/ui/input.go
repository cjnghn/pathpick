package ui

import (
	"fmt"
	"os"

	"github.com/atotto/clipboard"
	"github.com/cjnghn/pathpick/internal/tree"
	"github.com/eiannone/keyboard"
)

func (d *Display) eventLoop() error {
	if err := keyboard.Open(); err != nil {
		return err
	}
	defer keyboard.Close()

	for {
		d.render()

		char, key, err := keyboard.GetKey()
		if err != nil {
			return err
		}

		switch {
		case key == keyboard.KeyEsc || char == 'q':
			return nil
		case key == keyboard.KeyArrowUp:
			d.moveUp()
		case key == keyboard.KeyArrowDown:
			d.moveDown()
		case key == keyboard.KeyArrowLeft:
			d.moveOut()
		case key == keyboard.KeyArrowRight:
			d.moveIn()
		case key == keyboard.KeySpace:
			d.current.ToggleSelect()
		case char == 'y':
			return d.copySelected()
		}
	}
}

// moveDown moves the cursor to the next visible node
func (d *Display) moveDown() {
	var nextNode *tree.Node
	var allNodes []*tree.Node
	d.collectVisibleNodes(d.root, &allNodes)

	// Find current node in the list and move to next
	for i, node := range allNodes {
		if node == d.current && i < len(allNodes)-1 {
			nextNode = allNodes[i+1]
			break
		}
	}

	if nextNode != nil {
		d.current = nextNode
	}
}

// moveUp moves the cursor to the previous visible node
func (d *Display) moveUp() {
	var prevNode *tree.Node
	var allNodes []*tree.Node
	d.collectVisibleNodes(d.root, &allNodes)

	// Find current node in the list and move to previous
	for i, node := range allNodes {
		if node == d.current && i > 0 {
			prevNode = allNodes[i-1]
			break
		}
	}

	if prevNode != nil {
		d.current = prevNode
	}
}

// collectVisibleNodes collects all visible nodes in depth-first order
func (d *Display) collectVisibleNodes(node *tree.Node, nodes *[]*tree.Node) {
	*nodes = append(*nodes, node)
	for _, child := range node.Children {
		d.collectVisibleNodes(child, nodes)
	}
}

// moveIn changes focus to first child if exists
func (d *Display) moveIn() {
	if len(d.current.Children) > 0 {
		d.current = d.current.Children[0]
	}
}

// moveOut changes focus to parent if exists
func (d *Display) moveOut() {
	if d.current.Parent != nil {
		d.current = d.current.Parent
	}
}

func (d *Display) copySelected() error {
	selected := d.root.GetSelected()
	if len(selected) == 0 {
		return nil
	}

	var content string
	for _, node := range selected {
		data, err := os.ReadFile(node.Path)
		if err != nil {
			return err
		}
		content += fmt.Sprintf("=== %s ===\n", node.Path)
		content += string(data) + "\n\n"
	}

	return clipboard.WriteAll(content)
}
