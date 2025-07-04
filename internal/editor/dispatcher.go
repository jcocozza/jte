package editor

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/jcocozza/jte/internal/keyboard"
	"github.com/jcocozza/jte/internal/mode"
)

// the Dispatcher worries about correctly grouping keys together
//
// when it has a good grouping, it *dispatches* the grouping and associated actions to the main process
// for next steps
type Dispatcher struct {
	logger         *slog.Logger
	currKeys       keyboard.OrderedKeyList
	repeatModifier int
}

func NewDispatcher(l *slog.Logger) *Dispatcher {
	return &Dispatcher{
		logger:   l.WithGroup("dispatcher"),
		currKeys: keyboard.OrderedKeyList{},
	}
}

func (d *Dispatcher) accept(k keyboard.Key) {
	d.currKeys.Append(k)
	d.logger.Debug("accepted key", slog.String("key", k.String()), slog.String("current keys", d.currKeys.Collapse()), slog.Int("repeat modifier", d.repeatModifier))
}

// in insert mode:
// 1. check for a valid sequence (very few)
// 2. if valid, generate the action/changed based on that
// 3. otherwise, generate an insert action for next text
//
// return true to flush, false to continue
func (d *Dispatcher) processInsert(k keyboard.Key, n *BindingNode) (bool, []Action) {
	d.accept(k)
	possiblyValid := n.HasPrefix(d.currKeys)
	if possiblyValid {
		actionNode, err := n.Lookup(d.currKeys)
		if err != nil {
			return false, nil
		}
		return true, actionNode.Actions
	}
	return true, []Action{Insert{rune(d.currKeys[0])}}
}

// in command mode:
// 1. check for valid sequence (e.g. <enter>, <esc>, etc)
// 2. if valid dispatch command
// 3. otherwise, keep adding characters to command prompt
//
// return true to flush, false to continue
func (d *Dispatcher) processCommand(k keyboard.Key, n *BindingNode) (bool, []Action) {
	d.accept(k)
	possiblyValid := n.HasPrefix(d.currKeys)
	if possiblyValid {
		actionNode, err := n.Lookup(d.currKeys)
		if err != nil {
			return false, nil
		}
		return true, actionNode.Actions
	}
	return true, []Action{InsertCommandChar{rune(d.currKeys[0])}}
}

// in normal mode:
// 1. check for (possibly) valid sequence
// 2. if valid or possibly valid, keep appending until we get a valid or invalid
//
// return true to flush, false to continue
func (d *Dispatcher) processNormal(k keyboard.Key, n *BindingNode) (bool, []Action) {
	if k.IsDigit() {
		digit := int(k - '0')
		if d.repeatModifier == 0 {
			d.repeatModifier = digit
			return false, nil
		}
		d.repeatModifier = (d.repeatModifier * 10) + digit
		return false, nil
	}
	d.accept(k)
	possiblyValid := n.HasPrefix(d.currKeys)
	if possiblyValid {
		actionNode, err := n.Lookup(d.currKeys)
		if err != nil {
			return false, nil
		}
		if d.repeatModifier == 0 || d.repeatModifier == 1 {
			return true, actionNode.Actions
		}
		repeatedActions := []Action{} // TODO: allocate this properly
		for range d.repeatModifier {
			repeatedActions = append(repeatedActions, actionNode.Actions...)
		}
		return true, repeatedActions
	}
	return true, nil // since nothing matches, we just want to flush right away
}

var ErrNoDispatch = errors.New("no dispatch")

// this is run one time per event loop
//
// based on the mode, process the keypress accordingly
//
// will return ErrNoDispatch if nothing to report
func (d *Dispatcher) ProcessKeypress(k keyboard.Key, m mode.Mode, n *BindingNode) ([]Action, error) {

	var flush bool
	var actions []Action
	switch m {
	case mode.Normal:
		flush, actions = d.processNormal(k, n)
	case mode.Command:
		flush, actions = d.processCommand(k, n)
	case mode.Insert:
		flush, actions = d.processInsert(k, n)
	default:
		panic(fmt.Sprintf("invalid mode on dispatch: %s", m))
	}

	if !flush {
		return nil, ErrNoDispatch
	}

	d.repeatModifier = 0
	d.currKeys = keyboard.OrderedKeyList{}
	return actions, nil
}
