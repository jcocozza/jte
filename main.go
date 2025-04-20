package main

import (
	"github.com/jcocozza/jte/pkg/buffer"
	"github.com/jcocozza/jte/pkg/editor"
	"github.com/jcocozza/jte/pkg/logger"
	"github.com/jcocozza/jte/pkg/renderer"
)

func main() {
	l, f, err := logger.NewLogger()
	if err != nil {
		panic(err)
	}
	defer f.Close()

	e := editor.NewEditor(l)
	r := renderer.NewTextRenderer(l)
	err = r.Setup()
	if err != nil {
		panic("unable to setup")
	}

	id := e.BM.Add(*buffer.NewBuffer("[No Name]", buffer.SampleRows))
	e.BM.SetCurrent(id)
	for {
		exit, err := e.EventLoopStep()
		if err != nil {
			r.ExitErr(err)
		}
		if exit {
			r.Exit("")
		}
		r.Render(e)
	}
}
