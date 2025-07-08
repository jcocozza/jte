package main

import (
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
		panic(err)
	}

	buf, err := buffer.ReadFileIntoBuffer("/home/water/projects/jte/test.txt", l)
	if err != nil {
		panic(err)
	}

	e.BM.SetCurrent(e.BM.Add(buf))
	e.PM.Root.Bn = e.BM.Current

	r.Render(e)
	// event loop
	for {
		key, err := kb.GetKeypress()
		if err != nil {
			panic(err)
		}

		var n *action.BindingNode
		state := e.M.Current()

		switch state {
		case mode.Command:
		case mode.Insert:
		case mode.Normal:
		default:
			panic("invalid state")
		}

		actions, done := ap.AcceptKey(key, state, n)
		if !done {
			continue
		}

		for _, a := range actions {
			err := a.Apply(e)
			if err != nil {
				r.ExitErr(err)
			}
		}
		r.Render(e)
	}
}
