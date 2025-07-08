package editor

import (
	"log/slog"

	"github.com/jcocozza/jte/internal/buffer"
	"github.com/jcocozza/jte/internal/commmand"
	"github.com/jcocozza/jte/internal/mode"
	"github.com/jcocozza/jte/internal/panemanager"
)

type Editor struct {
	M *mode.ModeMachine
	BM *buffer.BufferManager
	PM *panemanager.PaneManager
	CW *commmand.CommandWindow
}

func NewEditor(l *slog.Logger) *Editor {
	return &Editor{
		M: mode.NewModeMachine(l),
		BM: buffer.NewBufferManager(l),
		PM: panemanager.NewPaneManager(l),
		CW: commmand.NewCommandWindow(l),
	}
}

func (e *Editor) Close() {
	e.PM.Delete()
	e.BM.Current = e.PM.Curr.Bn
}

func (e *Editor) Up() {
	e.PM.Up()
	e.BM.Current = e.PM.Curr.Bn
}
func (e *Editor) Down() {
	e.PM.Down()
	e.BM.Current = e.PM.Curr.Bn
}
func (e *Editor) Left() {
	e.PM.Left()
	e.BM.Current = e.PM.Curr.Bn
}
func (e *Editor) Right() {
	e.PM.Right()
	e.BM.Current = e.PM.Curr.Bn
}

// useful bits about the editor
type EditorStatus struct {
	Mode mode.Mode
	CurrentPane *panemanager.PaneNode
}

func (e *Editor) Status() *EditorStatus {
	return &EditorStatus{
		Mode: e.M.Current(),
		CurrentPane: e.PM.Curr,
	}
}
