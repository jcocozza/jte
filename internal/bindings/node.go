package bindings

import (
	"fmt"

	"github.com/jcocozza/jte/internal/keyboard"
)

type BindingNode struct {
	children map[keyboard.Key]*BindingNode
	Actions  []ActionId
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
		return n, nil
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
		keyboard.ESC: {
			Actions: []ActionId{Action_Commit, Action_ModeNormal},
		},
	},
}

var NormalBindings = &BindingNode{
	Actions: nil,
	children: map[keyboard.Key]*BindingNode{
		'i': { children: nil, Actions: []ActionId{Action_ModeInsert} },
	},
}

var CommandBindings = &BindingNode{
	Actions: nil,
	children: map[keyboard.Key]*BindingNode{
		keyboard.ESC: { children: nil, Actions: []ActionId{Action_ModeNormal} },
	},
}
