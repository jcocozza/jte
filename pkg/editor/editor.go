package editor

import (
	"fmt"
	"log/slog"

	"github.com/jcocozza/jte/pkg/actions"
	"github.com/jcocozza/jte/pkg/buffer"
	"github.com/jcocozza/jte/pkg/keyboard"
	"github.com/jcocozza/jte/pkg/state"
)

type Editor struct {
	kb *keyboard.Keyboard
	kq *keyboard.KeyQueue
	sm *state.StateMachine
	aq *actions.ActionQueue
	BM *buffer.BufferManager
}

func NewEditor(l *slog.Logger) *Editor {
	kb := keyboard.NewKeyboard(l)
	kq := keyboard.NewKeyQueue(l)
	sm := state.NewStateMachine()
	aq := actions.NewActionQueue(l)
	BM := buffer.NewBufferManager(l)
	return &Editor{
		kb: kb,
		kq: kq,
		sm: sm,
		aq: aq,
		BM: BM,
	}
}

// called in the main event loop
//
// return true if it is time to exit
func (e *Editor) EventLoopStep() (bool, error) {
	kp, err := e.kb.GetKeypress()
	if err != nil {
		return false, err
	}
	e.kq.Enqueue(kp)
	action := e.sm.HandleKeyQueue(e.kq)
	e.aq.Enqueue(action)
	for {
		action, err := e.aq.Dequeue()
		if err != nil { // end of queue
			break
		}
		if action == actions.Exit {
			return true, nil
		}
		fn, ok := Registry[action]
		if !ok { // a non existent action
			panic(fmt.Sprintf("action: %d does not exist in registry", action))
		}
		fn(e) //actually exectute the action
	}
	return false, nil
}

