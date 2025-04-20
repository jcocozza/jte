package bindings

import (
	"github.com/jcocozza/jte/pkg/actions"
	"github.com/jcocozza/jte/pkg/keyboard"
)

// this is a way to map consecutive key presses to different actions
//
// the 'root' will always have a None action
type BindingNode struct {
	children map[keyboard.Key]*BindingNode
	action   actions.Action
}

func RootBindingNode() *BindingNode {
	return &BindingNode{
		children: make(map[keyboard.Key]*BindingNode),
		action:   actions.None,
	}
}

// traverse the nodes in order of the queue
//
// will return the None action if no action is found
func (bn *BindingNode) Traverse(kq *keyboard.KeyQueue) actions.Action {
	curr := bn
	for {
		key, err := kq.Dequeue()
		if err != nil {
			// we are out of keys, time to break
			break
		}
		next, ok := curr.children[key]
		// break when next key isn't in the bindings
		if !ok {
			break
		}
		// this is necessary because we might write an
		// incomplete set of bindings
		if next == nil {
			panic("a keynode cannot have a child be nil when the key exists in the set of bindings")
		}
		curr = next
	}
	return curr.action
}
