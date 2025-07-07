package editor

import (
	"log/slog"

	"github.com/jcocozza/jte/internal/buffer"
	"github.com/jcocozza/jte/internal/mode"
	"github.com/jcocozza/jte/internal/panemanager"
)

type Editor struct {
	M *mode.ModeMachine
	BM *buffer.BufferManager
	PM *panemanager.PaneManager
}

func NewEditor(l *slog.Logger) *Editor {
	return &Editor{
		M: mode.NewModeMachine(l),
		BM: buffer.NewBufferManager(l),
		PM: panemanager.NewPaneManager(l),
	}
}
