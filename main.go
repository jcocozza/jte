package main

import (
	"github.com/jcocozza/jte/internal/buffer"
	"github.com/jcocozza/jte/internal/editor"
	"github.com/jcocozza/jte/internal/logger"
	"github.com/jcocozza/jte/internal/renderer"
)

func main() {
	f, err := logger.Init()
	if err != nil {
		panic(err)
	}
	defer f.Close()

	r := renderer.NewTextRenderer(logger.Logger)
	e := editor.NewEditor(logger.Logger)

	err = r.Setup()
	if err != nil {
		panic(err)
	}

	buf := buffer.NewBuffer("[No Name]", "", false, []buffer.BufRow{{'f', 'o', 'o'}}, logger.Logger)
	id := e.BM.Add(buf)
	e.BM.SetCurrent(id)

	// this is some setup stuff we're going to have to figure out
	p := &editor.Pane{nil, buf, true}
	n := &editor.SplitNode{Pane: p}
	e.Root = n
	e.Active = n
	r.Render(e) // initial render
	for {
		err := e.HandleKeypress()
		if err != nil {
			r.ExitErr(err)
		}
		r.Render(e)
	}
}
