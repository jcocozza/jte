package main

import (
	"log/slog"

	"github.com/jcocozza/jte/editor"
	"github.com/jcocozza/jte/logger"
)

func main() {
	l := logger.CreateLogger(slog.LevelDebug)
	e := editor.NewTextEditor(l)
	e.Run()
}
