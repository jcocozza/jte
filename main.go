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
	if err != nil {panic(err)}


	buf := buffer.NewBuffer("[No Name]", "", false, []buffer.BufRow{{'f','o','o'}}, logger.Logger)
	id := e.BM.Add(buf)
	e.BM.SetCurrent(id)

	p := &editor.Pane{nil, buf}
	e.SN = &editor.SplitNode{Pane: p}
	r.Render(e) // initial render
	for {
		err := e.HandleKeypress()
		if err != nil {
			r.ExitErr(err)
		}
		r.Render(e)
	}
}
