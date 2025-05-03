package state

import (
	"log/slog"

	"github.com/jcocozza/jte/pkg/actions"
	"github.com/jcocozza/jte/pkg/bindings"
	"github.com/jcocozza/jte/pkg/keyboard"
)

type InsertMode struct {
	bindings *bindings.BindingNode
	logger *slog.Logger
}

func NewInsertMode(l *slog.Logger) *InsertMode {
	b := bindings.RootBindingNode()
	return &InsertMode{
		bindings: b,
		logger: l.WithGroup("insert-mode"),
	}
}

func (m *InsertMode) IsPossiblyValid(kq []keyboard.Key) bool {
	if len(kq) == 1 && kq[0].IsUnicode() {
		return true
	}
	if m.bindings != nil {
		if !m.bindings.IsPossiblyValid(kq) {
			return bindings.Insert.IsPossiblyValid(kq)
		}
		return true
	}
	return bindings.Insert.IsPossiblyValid(kq)
}

func (m *InsertMode) Valid(kq []keyboard.Key) bool {
	if len(kq) == 1 && kq[0].IsUnicode() {
		return true
	}
	if m.bindings != nil {
		if !m.bindings.IsValid(kq) {
			return bindings.Insert.IsValid(kq)
		}
		return true
	}
	return bindings.Insert.IsValid(kq)
}

func (m *InsertMode) Name() ModeName {
	return Insert
}

func (m *InsertMode) HandleInput(kq *keyboard.KeyQueue) []actions.Action {
	duplicate := *kq
	newkq := &duplicate
	key, err := kq.Dequeue()
	if err != nil {
		return []actions.Action{}
	}
	if key.IsUnicode() {
		return []actions.Action{actions.InsertChar}
	}
	if m.bindings != nil {
		// use custom bindings first
		actionsList := m.bindings.Traverse(kq)
		if len(actionsList) > 0 {
			return actionsList
		}
		// try to use default bindings
		return bindings.Insert.Traverse(newkq)
	}
	return bindings.Insert.Traverse(kq)
}
