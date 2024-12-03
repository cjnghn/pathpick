package ui

import (
	"fmt"
	"os"

	"github.com/atotto/clipboard"
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

func (d *Display) moveUp() {
	if d.current.Parent == nil {
		return
	}
	siblings := d.current.Parent.Children
	for i, node := range siblings {
		if node == d.current && i > 0 {
			d.current = siblings[i-1]
			return
		}
	}
}

func (d *Display) moveDown() {
	if d.current.Parent == nil {
		if len(d.current.Children) > 0 {
			d.current = d.current.Children[0]
		}
		return
	}
	siblings := d.current.Parent.Children
	for i, node := range siblings {
		if node == d.current && i < len(siblings)-1 {
			d.current = siblings[i+1]
			return
		}
	}
}

func (d *Display) moveIn() {
	if len(d.current.Children) > 0 {
		d.current = d.current.Children[0]
	}
}

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
