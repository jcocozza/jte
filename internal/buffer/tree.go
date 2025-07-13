package buffer

import (
	"fmt"
	"log/slog"
)

type ChangeTracker struct {
	logger *slog.Logger
	Root   *ChangeNode
	Head   *ChangeNode
}

func NewChangeTracker(l *slog.Logger) *ChangeTracker {
	r := &ChangeNode{}
	return &ChangeTracker{
		logger: l.WithGroup("change-tracker"),
		Root:   r,
		Head:   r,
	}
}

func (ct *ChangeTracker) Redo(b *Buffer) error {
	if len(ct.Head.Children) == 0 {
		return fmt.Errorf("nothing to redo")
	}
	ct.Head.Cb.Do(b)
	ct.Head = ct.Head.Children[ct.Head.Active]
	return nil
}

func (ct *ChangeTracker) Undo(b *Buffer) error {
	if ct.Head.Parent == nil {
		return fmt.Errorf("nothing to undo")
	}
	ct.Head = ct.Head.Parent
	for _, c := range ct.Head.Cb.block {
		ct.logger.Debug(c.String())
	}
	ct.Head.Cb.Undo(b)
	return nil
}

func (ct *ChangeTracker) Add(c Change) {
	ct.logger.Debug("add", slog.String("change", c.String()))
	ct.Head.Cb.Add(c)
}

func (ct *ChangeTracker) Commit() {
	// don't commit if we don't have any changes
	if len(ct.Head.Cb.block) == 0 {
		ct.logger.Debug("nothing to commit")
		return
	}
	ct.Head.Cb.Lock()
	ct.logger.Debug("commit", slog.Any("block", ct.Head.Cb))
	n := &ChangeNode{Parent: ct.Head}
	ct.Head.Children = append(ct.Head.Children, n)
	ct.Head.Active = len(ct.Head.Children) - 1
	ct.Head = n
	ct.logger.Debug("commit", slog.String("diagram", ct.Root.draw("", false)))
}

// a doubly linked tree structure
type ChangeNode struct {
	Parent   *ChangeNode
	Children []*ChangeNode
	Active   int // idx of the latest child
	Cb       ChangeBlock
}

func (cn *ChangeNode) draw(prefix string, isLast bool) string {
	s := prefix
	if isLast {
		//s += "└─"
		s += "-"
	} else {
		s += "|-"
		//s += "-"
	}
	s += "*\n"

	for i, child := range cn.Children {
		isLastChild := i == len(cn.Children)-1
		newPrefix := prefix
		if isLast {
			newPrefix += "  "
		} else {
			//newPrefix += "│ "
			newPrefix += "| "
		}
		s += child.draw(newPrefix, isLastChild)
	}
	return s
}

//func (cn *ChangeNode) draw(spaces int) string {
//	s := strings.Repeat(" ", spaces) + "*" + "\n"
//	for _, child := range cn.Children {
//		s += child.draw(spaces+4)
//	}
//	return s
//}

func (cn *ChangeNode) String() string {
	return cn.draw("", false)
}
