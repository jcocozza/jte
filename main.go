package main

import (
	"log/slog"

	"github.com/jcocozza/jte/api/keypress"
	"github.com/jcocozza/jte/logger"
	"github.com/jcocozza/jte/term"
)

func main() {
	rw, err := term.EnableRawMode()
	if err != nil {
		panic(err)
	}
	defer rw.Restore()
	l := logger.CreateLogger(slog.LevelDebug)
	keyboard := keypress.NewKeyboard(l)
	for {
		_, err := keyboard.GetKeypress()
		if err != nil {
			panic(err)
		}
	}
}
