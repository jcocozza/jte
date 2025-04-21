package editor

import (
	"fmt"
	"log/slog"

	"github.com/jcocozza/jte/pkg/actions"
	"github.com/jcocozza/jte/pkg/buffer"
	commandwindow "github.com/jcocozza/jte/pkg/commandWindow"
	"github.com/jcocozza/jte/pkg/keyboard"
	"github.com/jcocozza/jte/pkg/state"
)

type Editor struct {
	kb *keyboard.Keyboard
	kq *keyboard.KeyQueue
	SM *state.StateMachine
	aq *actions.ActionQueue
	BM *buffer.BufferManager
	CW *commandwindow.CommandWindow

	logger *slog.Logger
}

func NewEditor(l *slog.Logger) *Editor {
	kb := keyboard.NewKeyboard(l)
	kq := keyboard.NewKeyQueue(l)
	sm := state.NewStateMachine(l)
	aq := actions.NewActionQueue(l)
	BM := buffer.NewBufferManager(l)
	CW := commandwindow.NewCommandWindow(l)
	return &Editor{
		kb: kb,
		kq: kq,
		SM: sm,
		aq: aq,
		BM: BM,
		CW: CW,
		logger: l.WithGroup("editor"),
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
	actionList := e.SM.HandleKeyQueue(e.kq)
	e.aq.EnqueueList(actionList)
	for { // for each action in the queue, do the action
		action, err := e.aq.Dequeue()
		if err != nil { // end of queue
			break
		}
		if action == actions.Exit {
			return true, nil
		}
		if action == actions.InsertChar {
			e.BM.Current.Buf.InsertChar(byte(kp))
			return false, nil
		}
		fn, ok := Registry[action]
		if !ok { // a non existent action
			panic(fmt.Sprintf("action: %d does not exist in registry", action))
		}
		fn(e) //actually exectute the action
	}
	return false, nil
}

