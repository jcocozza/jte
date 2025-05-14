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
		keyboard.ESC:   {Actions: []Action{SwitchMode{m: mode.Normal}}},
		keyboard.CtrlC: {children: nil, Actions: []Action{Exit{}}},

		keyboard.BACKSPACE:   {children: nil, Actions: []Action{Backspace{}}},
		keyboard.BACKSPACE_2: {children: nil, Actions: []Action{Backspace{}}},
		keyboard.DELETE:      {children: nil, Actions: []Action{Delete{}}},

		keyboard.TAB: {children: nil, Actions: []Action{Insert{c: rune(keyboard.TAB)}}},
		keyboard.ENTER: {children: nil, Actions: []Action{EnterNewLine{}}},

		keyboard.ARROW_UP:    {children: nil, Actions: []Action{CursorUp{}}},
		keyboard.ARROW_DOWN:  {children: nil, Actions: []Action{CursorDown{}}},
		keyboard.ARROW_LEFT:  {children: nil, Actions: []Action{CursorLeft{}}},
		keyboard.ARROW_RIGHT: {children: nil, Actions: []Action{CursorRight{}}},
	},
}

var NormalBindings = &BindingNode{
	Actions: nil,
	children: map[keyboard.Key]*BindingNode{
		'i':            {children: nil, Actions: []Action{SwitchMode{m: mode.Insert}}},
		keyboard.CtrlC: {children: nil, Actions: []Action{Exit{}}},

		'o': {children: nil, Actions: []Action{SwitchMode{m: mode.Insert}, NewLineBelow{}}},
		'O': {children: nil, Actions: []Action{SwitchMode{m: mode.Insert}, NewLineAbove{}}},

		'd': {Actions: nil,
			children: map[keyboard.Key]*BindingNode{
				'd': {children: nil, Actions: []Action{DeleteLine{}}},
			},
		},
		's': {children: nil, Actions: []Action{SplitHorizontal{}}},
		'v': {children: nil, Actions: []Action{SplitVertical{}}},

		'k': {children: nil, Actions: []Action{CursorUp{}}},
		'j': {children: nil, Actions: []Action{CursorDown{}}},
		'h': {children: nil, Actions: []Action{CursorLeft{}}},
		'l': {children: nil, Actions: []Action{CursorRight{}}},
	},
}

var CommandBindings = &BindingNode{
	Actions: nil,
	children: map[keyboard.Key]*BindingNode{
		keyboard.ESC:   {children: nil, Actions: []Action{SwitchMode{m: mode.Normal}}},
		keyboard.CtrlC: {children: nil, Actions: []Action{Exit{}}},
	},
}
