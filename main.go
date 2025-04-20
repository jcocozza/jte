package main

import (
	"fmt"
	"os"

	"github.com/jcocozza/jte/pkg/buffer"
	"github.com/jcocozza/jte/pkg/editor"
	"github.com/jcocozza/jte/pkg/fileutil"
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

	var buf *buffer.Buffer


	if len(os.Args) > 1 {
		path := os.Args[1]
		content, writeable, err := fileutil.ReadFile(path)
		if err != nil {
			r.ExitErr(fmt.Errorf("unable to read filepath: %w", err))
			return
		}
		readOnly := !writeable
		buf = buffer.NewBuffer(path, readOnly, content)
	} else {
		buf = buffer.NewBuffer("[No Name]", true, buffer.SampleRows)
	}

	id := e.BM.Add(buf)
	e.BM.SetCurrent(id)
	r.Render(e) // need to do an inital render
	// enter into the main event loop
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
