package panemanager

import (
	"log/slog"
)

// drawn from the root
//
// curr is a pointer keeping track of the current pane
type PaneManager struct {
	logger *slog.Logger
	root *PaneNode
	curr *PaneNode
}

func (p *PaneManager) Left() {
	p.logger.Debug("move left")
	p.curr = p.curr.Left()
}
func (p *PaneManager) Right() {
	p.logger.Debug("move right")
	p.curr = p.curr.Right()
}
func (p *PaneManager) Up() {
	p.logger.Debug("move up")
	p.curr = p.curr.Up()
}
func (p *PaneManager) Down() {
	p.logger.Debug("move down")
	p.curr = p.curr.Down()
}

func (p *PaneManager) Delete() {
	p.logger.Debug("delete pane")
	curr, root := p.curr.Delete()
	if root != nil {
		p.root = root
	}
	p.curr = curr
}
