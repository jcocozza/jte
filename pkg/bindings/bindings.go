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
	actions  []actions.Action
}

func RootBindingNode() *BindingNode {
	return &BindingNode{
		children: make(map[keyboard.Key]*BindingNode),
		actions:  []actions.Action{},
	}
}

// traverse the nodes in order of the queue
//
// will return the None action if no action is found
func (bn *BindingNode) Traverse(kq *keyboard.KeyQueue) []actions.Action {
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
	return curr.actions
}

func (bn *BindingNode) IsPossiblyValid(kq []keyboard.Key) bool {
	curr := bn
	for _, key := range kq {
		next, ok := curr.children[key]
		if !ok {
			return false
		}
		curr = next
	}
	return true
}

// a queue of keys is valid if it returns one or more actions
func (bn *BindingNode) IsValid(kq []keyboard.Key) bool {
	curr := bn
	for _, key := range kq {
		next, ok := curr.children[key]
		if !ok {
			break
		}
		if next == nil {
			panic("a keynode cannot have a child be nil when the key exists in the set of bindings")
		}
		curr = next
	}
	return len(curr.actions) > 0
}
