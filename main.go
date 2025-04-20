package main

import (
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
	e.Run()
}
