package dispatcher

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/jcocozza/jte/internal/bindings"
	"github.com/jcocozza/jte/internal/keyboard"
	"github.com/jcocozza/jte/internal/mode"
)

type Dispatch struct {
	Keys    keyboard.OrderedKeyList
	Actions []bindings.ActionId
}

// the Dispatcher worries about correctly grouping keys together
//
// when it has a good grouping, it *dispatches* the grouping and associated actions to the main process
// for next steps
type Dispatcher struct {
	logger   *slog.Logger
	currKeys keyboard.OrderedKeyList
}

func NewDispatcher(l *slog.Logger) *Dispatcher {
	return &Dispatcher{
		logger:   l.WithGroup("dispatcher"),
		currKeys: keyboard.OrderedKeyList{},
	}
}

func (d *Dispatcher) accept(k keyboard.Key) {
	d.logger.Debug("accepting key", slog.String("key", k.String()))
	d.currKeys.Append(k)
}

// in insert mode:
// 1. check for a valid sequence (very few)
// 2. if valid, generate the action/changed based on that
// 3. otherwise, generate an insert action for next text
//
// return true to flush, false to continue
func (d *Dispatcher) processInsert(n *bindings.BindingNode) (bool, []bindings.ActionId) {
	possiblyValid := n.HasPrefix(d.currKeys)
	if possiblyValid {
		actionNode, err := n.Lookup(d.currKeys)
		if err != nil {
			return false, nil
		}
		return true, actionNode.Actions
	}
	return true, []bindings.ActionId{bindings.Action_InsertChar}
}

// in command mode:
// 1. check for valid sequence (e.g. <enter>, <esc>, etc)
// 2. if valid dispatch command
// 3. otherwise, keep adding characters to command prompt
//
// return true to flush, false to continue
func (d *Dispatcher) processCommand(n *bindings.BindingNode) (bool, []bindings.ActionId) {
	possiblyValid := n.HasPrefix(d.currKeys)
	if possiblyValid {
		actionNode, err := n.Lookup(d.currKeys)
		if err != nil {
			return false, nil
		}
		return true, actionNode.Actions
	}
	return true, []bindings.ActionId{bindings.InsertCommandChar}
}

// in normal mode:
// 1. check for (possibly) valid sequence
// 2. if valid or possibly valid, keep appending until we get a valid or invalid
//
// return true to flush, false to continue
func (d *Dispatcher) processNormal(n *bindings.BindingNode) (bool, []bindings.ActionId) {
	possiblyValid := n.HasPrefix(d.currKeys)
	if possiblyValid {
		actionNode, err := n.Lookup(d.currKeys)
		if err != nil {
			return false, nil
		}
		return true, actionNode.Actions
	}
	return true, nil // since nothing matches, we just want to flush right away
}

var ErrNoDispatch = errors.New("no dispatch")

// this is run one time per event loop
//
// based on the mode, process the keypress accordingly
//
// will return ErrNoDispatch if nothing to report
func (d *Dispatcher) ProcessKeypress(k keyboard.Key, m mode.Mode, n *bindings.BindingNode) (Dispatch, error) {
	d.accept(k)

	var flush bool
	var actions []bindings.ActionId
	switch m {
	case mode.Normal:
		flush, actions = d.processNormal(n)
	case mode.Command:
		flush, actions = d.processCommand(n)
	case mode.Insert:
		flush, actions = d.processInsert(n)
	default:
		panic(fmt.Sprintf("invalid mode on dispatch: %s", m))
	}

	if !flush {
		return Dispatch{}, ErrNoDispatch
	}

	dispatch := Dispatch{
		Keys:    d.currKeys,
		Actions: actions,
	}
	d.currKeys = keyboard.OrderedKeyList{}
	return dispatch, nil
}
