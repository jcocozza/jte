package editor

import (
	"log/slog"

	"github.com/jcocozza/jte/pkg/actions"
	"github.com/jcocozza/jte/pkg/keyboard"
	"github.com/jcocozza/jte/pkg/state"
)

type Editor struct {
	kb *keyboard.Keyboard
	kq *keyboard.KeyQueue
	sm *state.StateMachine
	aq *actions.ActionQueue
}

func NewEditor(l *slog.Logger) *Editor {
	kb := keyboard.NewKeyboard(l)
	kq := keyboard.NewKeyQueue(l)
	sm := state.NewStateMachine()
	aq := actions.NewActionQueue(l)
	return &Editor{
		kb: kb,
		kq: kq,
		sm: sm,
		aq: aq,
	}
}

func (e *Editor) Run() {
	for {
		kp, err := e.kb.GetKeypress()
		if err != nil {
			panic(err)
		}
		e.kq.Enqueue(kp)
		action := e.sm.HandleKeyQueue(e.kq)
		if action != nil {
			e.aq.Enqueue(action)
		}
		// in the future, we may want to append other actions
		// for now the 'action queue' really on ever gets 1 action
		e.aq.Process()
	}
}
