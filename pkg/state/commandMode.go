package state

import (
	"github.com/jcocozza/jte/pkg/actions"
	"github.com/jcocozza/jte/pkg/bindings"
	"github.com/jcocozza/jte/pkg/keyboard"
)

type CommandMode struct {
	bindings *bindings.BindingNode
}

func NewCommandMode() *CommandMode {
	b := bindings.RootBindingNode()
	return &CommandMode{
		bindings: b,
	}
}

func (m *CommandMode) Name() ModeName {
	return Command
}

func (m *CommandMode) IsPossiblyValid(kq []keyboard.Key) bool {
	if m.bindings != nil {
		if !m.bindings.IsPossiblyValid(kq) {
			return bindings.Command.IsPossiblyValid(kq)
		}
		return true
	}
	return bindings.Command.IsPossiblyValid(kq)
}

func (m *CommandMode) Valid(kq []keyboard.Key) bool {
	if m.bindings != nil {
		if !m.bindings.IsValid(kq) {
			return bindings.Command.IsValid(kq)
		}
		return true
	}
	return bindings.Command.IsValid(kq)
}

func (m *CommandMode) HandleInput(kq *keyboard.KeyQueue) []actions.Action {
	duplicate := *kq
	newkq := &duplicate
	key, err := kq.Dequeue()
	if err != nil {
		return []actions.Action{}
	}
	if key.IsUnicode() {
		return []actions.Action{actions.InsertCommandChar}
	}
	if m.bindings != nil {
		// use custom bindings first
		actionsList := m.bindings.Traverse(kq)
		if len(actionsList) > 0 {
			return actionsList
		}
		// try to use default bindings
		return bindings.Command.Traverse(newkq)
	}
	return bindings.Command.Traverse(kq)
}
