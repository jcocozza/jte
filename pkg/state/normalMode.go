package state

import (
	"github.com/jcocozza/jte/pkg/actions"
	"github.com/jcocozza/jte/pkg/bindings"
	"github.com/jcocozza/jte/pkg/keyboard"
)

type NormalMode struct {
	bindings *bindings.BindingNode
}

func NewNormalMode() *NormalMode {
	b := bindings.RootBindingNode()
	return &NormalMode{
		bindings: b,
	}
}

func (m *NormalMode) Name() string {
	return "normal"
}

func (m *NormalMode) HandleInput(kq *keyboard.KeyQueue) actions.Action {
	// we need to duplicate the key queue incase we need to retry with default bindings
	// TODO: there has to be a better way to do this
	// we probably should just not dequeue things
	if m.bindings != nil {
		duplicate := *kq
		newkq := &duplicate
		// use custom bindings first
		action := m.bindings.Traverse(kq)
		if action != actions.None {
			return action
		}
		// try to use default bindings
		return bindings.Normal.Traverse(newkq)
	}
	return bindings.Normal.Traverse(kq)
}
