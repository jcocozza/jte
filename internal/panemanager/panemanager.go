package panemanager

import (
	"log/slog"
)

// drawn from the Root
//
// Curr is a pointer keeping track of the Current pane
type PaneManager struct {
	logger *slog.Logger
	Root   *PaneNode
	Curr   *PaneNode
}

// initialize pane manager with 1 Root pane
func NewPaneManager(l *slog.Logger) *PaneManager {
	r := newRootPaneNode()
	return &PaneManager{
		logger: l.WithGroup("pane-manager"),
		Root:   r,
		Curr:   r,
	}
}

func (p *PaneManager) Vsplit() {
	p.logger.Debug("v split")
	p.Curr = p.Curr.splitVertical()
}

func (p *PaneManager) Hsplit() {
	p.logger.Debug("h split")
	p.Curr = p.Curr.splitHorizontal()
}

func (p *PaneManager) Left() {
	p.logger.Debug("move left")
	p.Curr = p.Curr.Left()
}

func (p *PaneManager) Right() {
	p.logger.Debug("move right")
	p.Curr = p.Curr.Right()
}

func (p *PaneManager) Up() {
	p.logger.Debug("move up")
	p.Curr = p.Curr.Up()
}

func (p *PaneManager) Down() {
	p.logger.Debug("move down")
	p.Curr = p.Curr.Down()
}

func (p *PaneManager) Delete() {
	p.logger.Debug("delete pane")
	Curr, Root := p.Curr.Delete()
	if Root != nil {
		p.Root = Root
	}
	p.Curr = Curr
}
