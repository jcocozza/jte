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
