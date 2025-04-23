package state

import (
	"log/slog"
	"strconv"

	"github.com/jcocozza/jte/pkg/actions"
	"github.com/jcocozza/jte/pkg/bindings"
	"github.com/jcocozza/jte/pkg/keyboard"
)

type NormalMode struct {
	bindings *bindings.BindingNode
	logger *slog.Logger
}

func NewNormalMode(l *slog.Logger) *NormalMode {
	b := bindings.RootBindingNode()
	return &NormalMode{
		bindings: b,
		logger: l.WithGroup("normal-mode"),
	}
}

func (m *NormalMode) Name() ModeName {
	return Normal
}

func (m *NormalMode) IsPossiblyValid(kq []keyboard.Key) bool {
	i := 0
	for {
		if i > len(kq) - 1 {
			return true // the case were the user is still entering a number
		}
		if kq[i].IsDigit() {
			i++
		} else {
			break
		}
	}
	// after the user has entered a number, need to check to see if the following sequence is valid
	if m.bindings != nil {
		if m.bindings.IsPossiblyValid(kq[i:]) {
			return true
		} else {
			return bindings.Normal.IsPossiblyValid(kq[i:])
		}
	}
	return bindings.Normal.IsPossiblyValid(kq[i:])
}

func (m *NormalMode) Valid(kq []keyboard.Key) bool {
	i := 0
	for {
		if i > len(kq) - 1 {
			return false // just numbers is not a valid sequence
		}
		if kq[i].IsDigit() {
			i++
		} else {
			break
		}
	}
	if m.bindings != nil {
		if m.bindings.IsValid(kq[i:]) {
			return true
		} else {
			return bindings.Normal.IsValid(kq[i:])
		}
	}
	return bindings.Normal.IsValid(kq[i:])
}

func (m *NormalMode) HandleInput(kq *keyboard.KeyQueue) []actions.Action {
	numStr := ""
	dqCnt := 0
	for _, key := range kq.Keys() {
		if key.IsDigit() {
			numStr += strconv.Itoa(int(key - '0'))
			dqCnt++
		}
	}
	multipler, err := strconv.Atoi(numStr)
	if err != nil {
		multipler = 0
	}
	m.logger.Debug("multiplier", slog.Int("multiplier", multipler))
	for i := 0; i < dqCnt; i++ {
		kq.Dequeue()
	}
	// we need to duplicate the key queue incase we need to retry with default bindings
	// TODO: there has to be a better way to do this
	// we probably should just not dequeue things
	acts := []actions.Action{}
	if m.bindings != nil {
		duplicate := *kq
		newkq := &duplicate
		// use custom bindings first
		actionList := m.bindings.Traverse(kq)
		if len(actionList) > 0 {
			//return actionList
			acts = append(acts, actionList...)
		}
		// try to use default bindings
		//return bindings.Normal.Traverse(newkq)
		acts = append(acts, bindings.Normal.Traverse(newkq)...)
	}
	acts = append(acts, bindings.Normal.Traverse(kq)...)
	//return bindings.Normal.Traverse(kq)
	if multipler > 0 {
		var resultActs []actions.Action
		for i := 0; i < multipler; i++ {
			resultActs = append(resultActs, acts...)
		}
		return resultActs
	}
	return acts
}
