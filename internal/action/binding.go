package action

import (
	"fmt"

	"github.com/jcocozza/jte/internal/keyboard"
)

type BindingNode struct {
	children map[keyboard.Key]*BindingNode
	// the actions associated with the particular binding
	Actions []Action
}

// check if a list of keys is valid
func (n *BindingNode) HasPrefix(keys []keyboard.Key) bool {
	if len(keys) == 0 {
		return true
	}
	child, ok := n.children[keys[0]]
	if !ok {
		return false
	}
	return child.HasPrefix(keys[1:])
}

// return the node that is the result of traversing the bindings
func (n *BindingNode) Lookup(keys []keyboard.Key) (*BindingNode, error) {
	if len(keys) == 0 {
		if len(n.Actions) > 0 {
			return n, nil
		}
		return nil, fmt.Errorf("invalid key sequence, no leaf node: %s", keyboard.Collapse(keys))
	}
	child, ok := n.children[keys[0]]
	if !ok {
		return nil, fmt.Errorf("invalid key sequence %s", keyboard.Collapse(keys))
	}
	return child.Lookup(keys[1:])
}
