package panemanager

import (
	"log/slog"
)

type PaneManager struct {
	logger *slog.Logger
	root *PaneNode
	curr *PaneNode
}

func (p *PaneManager) Left() {
	p.curr = p.curr.Left()
}
func (p *PaneManager) Right() {
	p.curr = p.curr.Right()
}
func (p *PaneManager) Up() {
	p.curr = p.curr.Up()
}
func (p *PaneManager) Down() {
	p.curr = p.curr.Down()
}

func (p *PaneManager) Delete() {
	curr, root := p.curr.Delete()
	if root != nil {
		p.root = root
	}
	p.curr = curr
}
