package state

import (
	"fmt"
	"log/slog"

	"github.com/jcocozza/jte/pkg/actions"
	"github.com/jcocozza/jte/pkg/keyboard"
)

type Mode interface {
	Name() ModeName
	HandleInput(kq *keyboard.KeyQueue) []actions.Action
	IsPossiblyValid(kq []keyboard.Key) bool
	Valid(kq []keyboard.Key) bool
}

type StateMachine struct {
	current Mode
	modes   map[ModeName]Mode
	logger *slog.Logger
}

// setup the state machine
//
// currently this just sets up the different modes
func NewStateMachine(l *slog.Logger) *StateMachine {
	modes := make(map[ModeName]Mode)
	s := &StateMachine{
		modes: modes,
		logger: l.WithGroup("state-machine"),
	}

	normal := &NormalMode{}
	insert := NewInsertMode(l)
	command := NewCommandMode()
	s.register(normal)
	s.register(insert)
	s.register(command)

	s.SetMode(Normal)
	return s
}

func (sm *StateMachine) Current() string {
	return string(sm.current.Name())
}

func (sm *StateMachine) register(m Mode) {
	if _, ok := sm.modes[m.Name()]; ok {
		errMsg := fmt.Sprintf("mode: %s already registered", m.Name())
		panic(errMsg)
	}
	sm.modes[m.Name()] = m
}

func (sm *StateMachine) SetMode(name ModeName) {
	if m, ok := sm.modes[name]; ok {
	sm.logger.Debug("set mode", slog.String("mode", string(name)))
		sm.current = m
		return
	}
	panic("unexpected mode")
}

func (sm *StateMachine) IsPossiblyValid(keys []keyboard.Key) bool {
	if len(keys) == 0 {
		return false
	}
	return sm.current.IsPossiblyValid(keys)
}

func (sm *StateMachine) ValidSequence(keys []keyboard.Key) bool {
	if len(keys) == 0 {
		return false
	}
	return sm.current.Valid(keys)
}

func (sm *StateMachine) HandleKeyQueue(kq *keyboard.KeyQueue) []actions.Action {
	return sm.current.HandleInput(kq)
}
