package main

import (
	"os"

	"github.com/jcocozza/jte/pkg/editor"
	"github.com/jcocozza/jte/pkg/logger"
	"github.com/jcocozza/jte/pkg/term"
)

func main() {
	rw, err := term.EnableRawMode()
	if err != nil {
		panic(err)
	}
	defer rw.Restore()
	l, f, err := logger.NewLogger()
	if err != nil {
		panic(err)
	}
	defer f.Close()

	e := editor.NewEditor(l)
	for {
		err := e.EventLoopStep()
		if err != nil {
			rw.Restore()
			os.Exit(0)
		}
	}
}
