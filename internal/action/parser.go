package action

import (
	"log/slog"

	"github.com/jcocozza/jte/internal/keyboard"
	"github.com/jcocozza/jte/internal/mode"
)

// the action parser converts groups of keystrokes into actions
//
// when it has a complete set of keys, it "flushes" and returns a list of actions
type ActionParser struct {
	logger      *slog.Logger
	currentKeys []keyboard.Key
	repeatModifier int
}

func NewActionParser(l *slog.Logger) *ActionParser {
	return &ActionParser{
		logger:      l.WithGroup("action-parser"),
		currentKeys: []keyboard.Key{},
	}
}

func (ap *ActionParser) appendKey(key keyboard.Key) {
	ap.logger.Debug("append key", slog.String("key", key.String()))
	ap.currentKeys = append(ap.currentKeys, key)
}

func (ap *ActionParser) flush() {
	ap.logger.Debug("flush")
	ap.currentKeys = []keyboard.Key{}
}

// in normal mode:
// 1. check for (possibly) valid sequence
// 2. if valid or possibly valid, keep appending until we get a valid or invalid
//
// the bool will be true if we have completed parsing
func (ap *ActionParser) parseNormal(n *BindingNode) ([]Action, bool) {
	possiblyValid := n.HasPrefix(ap.currentKeys)
	if possiblyValid {
		actionNode, err := n.Lookup(ap.currentKeys)
		if err != nil {
			return nil, false
		}
		if ap.repeatModifier == 0 || ap.repeatModifier == 1 {
			return actionNode.Actions, true
		}
		repeatedActions := []Action{} // TODO: allocate this properly
		for range ap.repeatModifier {
			repeatedActions = append(repeatedActions, actionNode.Actions...)
		}
		return repeatedActions, true
	}
	return  nil, true // since nothing matches, we just want to flush right away
}

// the bool will be true if we have completed parsing
func (ap *ActionParser) parseInsert() ([]Action, bool) {
	return nil, false
}

// the bool will be true if we have completed parsing
func (ap *ActionParser) parseCommand() ([]Action, bool) {
	return nil, false
}

// this is run one time per event loop
//
// based on the mode, process the keypress accordingly
//
// return true if the full set of actions is ready to go.
func (ap *ActionParser) AcceptKey(key keyboard.Key, m mode.Mode, b *BindingNode) ([]Action, bool) {
	var actions []Action
	var done bool
	switch m {
	case mode.Normal:
		ap.appendKey(key)
		actions, done = ap.parseNormal(NormalBindings)
	case mode.Insert:
		ap.appendKey(key)
		actions, done = ap.parseInsert()
	case mode.Command:
		ap.appendKey(key)
		actions, done = ap.parseCommand()
	}
	if done {
		ap.flush()
	}
	return actions, done
}
