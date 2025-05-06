package editor

import (
	"log/slog"

	"github.com/jcocozza/jte/internal/buffer"
	"github.com/jcocozza/jte/internal/keyboard"
	"github.com/jcocozza/jte/internal/mode"
)

type Editor struct {
	kb *keyboard.Keyboard
	m  *mode.StateMachine
	d  *Dispatcher
	BM *buffer.BufferManager

	Root *SplitNode
	Active *SplitNode

	logger *slog.Logger
}

func NewEditor(l *slog.Logger) *Editor {
	return &Editor{
		kb: keyboard.NewKeyboard(l),
		m:  mode.NewStateMachine(l),
		d:  NewDispatcher(l),
		BM: buffer.NewBufferManager(l),
		Root: nil,
		Active: nil,

		logger: l.WithGroup("editor"),
	}
}

func (e *Editor) Mode() string {
	return string(e.m.Current())	
}

func (e *Editor) HandleKeypress() error {
	k, err := e.kb.GetKeypress()
	if err != nil {
		return err
	}
	var n *BindingNode
	state := e.m.Current()
	switch state {
	case mode.Command:
		n = CommandBindings
	case mode.Insert:
		n = InsertBindings
	case mode.Normal:
		n = NormalBindings
	default:
		panic("invalid state")
	}
	actions, err := e.d.ProcessKeypress(k, state, n)
	// no dispatch, nothing to do
	if err != nil {
		return nil
	}
	for _, action := range actions {
		e.logger.Debug("applying action", slog.String("action", action.String()))
		err := action.Apply(e)
		if err != nil {
			return err
		}
	}
	return nil
}
