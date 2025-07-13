package main

import (
	"errors"
	"fmt"
	"os"
	"runtime/debug"

	"github.com/jcocozza/jte/internal/action"
	"github.com/jcocozza/jte/internal/buffer"
	"github.com/jcocozza/jte/internal/editor"
	"github.com/jcocozza/jte/internal/keyboard"
	"github.com/jcocozza/jte/internal/logger"
	"github.com/jcocozza/jte/internal/mode"
	"github.com/jcocozza/jte/internal/renderer"
)

func main() {
	l := logger.NewLogger()
	kb := keyboard.NewKeyboard(l)
	ap := action.NewActionParser(l)
	r := renderer.NewRenderer(l)

	e := editor.NewEditor(l)

	err := r.Setup()
	if err != nil {
		r.ExitErr(err)
	}

	var buf *buffer.Buffer
	if len(os.Args) <= 1 {
		buf = buffer.NewEmptyBuffer(l)
	} else {
		buf, err = buffer.ReadFileIntoBuffer(os.Args[1], l)
		if err != nil {
			r.ExitErr(err)
		}
	}

	e.BM.SetCurrent(e.BM.Add(buf))
	e.PM.Root.Bn = e.BM.Current

	// this catches panics
	// and allows us to restore the terminal gracefully
	defer func() {
		if f := recover(); f != nil {
			stack := debug.Stack()
			r.ExitErr(fmt.Errorf("panic: %v\n\n%s", f, stack))
		}
	}()

	r.Render(e)
	// event loop
	for {
		key, err := kb.GetKeypress()
		if err != nil {
			r.ExitErr(err)
		}

		var n *action.BindingNode
		state := e.M.Current()

		switch state {
		case mode.Command:
		case mode.Insert:
		case mode.Normal:
		default:
			r.ExitErr(fmt.Errorf("invalid state"))
		}

		actions, done := ap.AcceptKey(key, state, n)
		if !done {
			continue
		}
		// this is an UGLY way to do this
		// we need to find a better way to reset the modifier
		if len(actions) > 0 {
			ap.ResetRepeat()
		}

		for _, a := range actions {
			err := a.Apply(e)
			switch {
			case errors.Is(err, action.ErrExit):
				r.ExitErr(err)
			case err != nil:
				e.CW.PushErr(err)
			default:
				continue
			}
		}
		r.Render(e)
	}
}
