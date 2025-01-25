package cli

import (
	"log/slog"
	"os"

	"github.com/jcocozza/jte/editor"
	"github.com/jcocozza/jte/logger"
)

func CLI() {
	l := logger.CreateLogger(slog.LevelDebug)
	e := editor.NewTextEditor(l)
	if len(os.Args) < 2 {
		e.Run()
	} else {
		fname := os.Args[1]
		err := e.Open(fname)
		if err != nil {
			panic(err)
		}
		e.Run()
	}
}
