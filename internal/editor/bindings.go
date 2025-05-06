package editor

import (
	"fmt"

	"github.com/jcocozza/jte/internal/keyboard"
	"github.com/jcocozza/jte/internal/mode"
)

type BindingNode struct {
	children map[keyboard.Key]*BindingNode
	Actions  []Action
}

func (n *BindingNode) IsLeaf() bool {
	return len(n.Actions) > 0
}

// bwhahahaha classic recursion
//
// check if the list of keys is valid
func (n *BindingNode) HasPrefix(keys keyboard.OrderedKeyList) bool {
	if len(keys) == 0 {
		return true
	}
	child, ok := n.children[keys[0]]
	if !ok {
		return false
	}
	return child.HasPrefix(keys[1:])
}

// forcing an error here because I will be lazy if I don't
//
// return the node that is the result of traversing the bindings
func (n *BindingNode) Lookup(keys keyboard.OrderedKeyList) (*BindingNode, error) {
	if len(keys) == 0 {
		if n.IsLeaf() {
			return n, nil
		}
		return nil, fmt.Errorf("invalid key sequence, no leaf node: %s", keys.Collapse())
	}
	child, ok := n.children[keys[0]]
	if !ok {
		return nil, fmt.Errorf("invalid key sequence %s", keys.Collapse())
	}
	return child.Lookup(keys[1:])
}

var InsertBindings = &BindingNode{
	Actions: nil,
	children: map[keyboard.Key]*BindingNode{
		keyboard.ESC:   {Actions: []Action{Commit{}, SwitchMode{m: mode.Normal}}},
		keyboard.CtrlC: {children: nil, Actions: []Action{Exit{}}},
	},
}

var NormalBindings = &BindingNode{
	Actions: nil,
	children: map[keyboard.Key]*BindingNode{
		'i':            {children: nil, Actions: []Action{SwitchMode{m: mode.Insert}}},
		keyboard.CtrlC: {children: nil, Actions: []Action{Exit{}}},

		'd': {Actions: nil,
			children: map[keyboard.Key]*BindingNode{
				'd': {children: nil, Actions: []Action{DeleteLine{}}},
			},
		},
	},
}

var CommandBindings = &BindingNode{
	Actions: nil,
	children: map[keyboard.Key]*BindingNode{
		keyboard.ESC:   {children: nil, Actions: []Action{SwitchMode{m: mode.Normal}}},
		keyboard.CtrlC: {children: nil, Actions: []Action{Exit{}}},
	},
}
