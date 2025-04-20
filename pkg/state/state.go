package state

import (
	"fmt"

	"github.com/jcocozza/jte/pkg/actions"
	"github.com/jcocozza/jte/pkg/keyboard"
)

type Mode interface {
	Name() string
	HandleInput(kq *keyboard.KeyQueue) actions.Action
}

type StateMachine struct {
	current Mode
	modes   map[string]Mode
}

// setup the state machine
//
// currently this just sets up the different modes
func NewStateMachine() *StateMachine {
	modes := make(map[string]Mode)
	s := &StateMachine{
		modes: modes,
	}

	normal := &NormalMode{}
	s.register(normal)
	s.SetMode("normal")
	return s
}

func (sm *StateMachine) register(m Mode) {
	if _, ok := sm.modes[m.Name()]; ok {
		errMsg := fmt.Sprintf("mode: %s already registered", m.Name())
		panic(errMsg)
	}
	sm.modes[m.Name()] = m
}

func (sm *StateMachine) SetMode(name string) {
	if m, ok := sm.modes[name]; ok {
		sm.current = m
		return
	}
	panic("unexpected mode")
}

func (sm *StateMachine) HandleKeyQueue(kq *keyboard.KeyQueue) actions.Action {
	return sm.current.HandleInput(kq)
}
