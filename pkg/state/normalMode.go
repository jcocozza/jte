package state

import (
	"github.com/jcocozza/jte/pkg/actions"
	"github.com/jcocozza/jte/pkg/keyboard"
)

type NormalMode struct {
	bindings *KeyNode
}

func (m *NormalMode) Name() string {
	return "normal"
}

// todo, finish this out
var defaultNormalBindings = &KeyNode{
	children: map[keyboard.Key]*KeyNode{
		'w': {children: nil, action: nil},
		'h': {children: nil, action: actions.CursorLeft},
		'j': {children: nil, action: actions.CursorDown},
		'k': {children: nil, action: actions.CursorUp},
		'l': {children: nil, action: actions.CursorRight},
		//keyboard.CtrlC: nil,
	},
	action: nil,
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
		if action != nil {
			return action
		}
		// try to use default bindings
		return defaultNormalBindings.Traverse(newkq)
	}
	return defaultNormalBindings.Traverse(kq)
}
