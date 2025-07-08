package action

import (
	"errors"
	"fmt"

	"github.com/jcocozza/jte/internal/editor"
	"github.com/jcocozza/jte/internal/keyboard"
	"github.com/jcocozza/jte/internal/mode"
)

// An Action is something done to the editor
type Action interface {
	// for debugging
	String() string
	Apply(e *editor.Editor) error
}

var ErrExit = errors.New("exit")

type Exit struct{}

func (a Exit) String() string { return "exit" }
func (a Exit) Apply(e *editor.Editor) error {
	return ErrExit
}

type SwitchMode struct {
	m mode.Mode
}

func (a SwitchMode) String() string { return fmt.Sprintf("switch mode: %d", a.m) }
func (a SwitchMode) Apply(e *editor.Editor) error {
	switch a.m {
	case mode.Insert:
	case mode.Normal:
	case mode.Command:
		e.CW.Unlock()
		e.CW.Hide()
	default:
		panic("nothing to do there")
	}
	e.M.SetMode(a.m)
	return nil
}

type CursorUp struct{}

func (a CursorUp) String() string               { return "CursorUp" }
func (a CursorUp) Apply(e *editor.Editor) error { e.BM.Current.Buf.Up(); return nil }

type CursorDown struct{}

func (a CursorDown) String() string               { return "CursorDown" }
func (a CursorDown) Apply(e *editor.Editor) error { e.BM.Current.Buf.Down(); return nil }

type CursorLeft struct{}

func (a CursorLeft) String() string               { return "CursorLeft" }
func (a CursorLeft) Apply(e *editor.Editor) error { e.BM.Current.Buf.Left(); return nil }

type CursorRight struct{}

func (a CursorRight) String() string               { return "CursorRight" }
func (a CursorRight) Apply(e *editor.Editor) error { e.BM.Current.Buf.Right(); return nil }

type SplitVertical struct{}

func (a SplitVertical) String() string { return "vert split" }
func (a SplitVertical) Apply(e *editor.Editor) error {
	e.PM.Vsplit()
	return nil
}

type SplitHorizontal struct{}

func (a SplitHorizontal) String() string { return "horizontal split" }
func (a SplitHorizontal) Apply(e *editor.Editor) error {
	e.PM.Hsplit()
	return nil
}

type SplitClose struct{}

func (a SplitClose) String() string { return "close split" }
func (a SplitClose) Apply(e *editor.Editor) error {
	e.Close()
	return nil
}

type PaneUp struct{}
func (a PaneUp) String() string { return "pane up" }
func (a PaneUp) Apply(e *editor.Editor) error {
	e.Up()
	return nil
}
type PaneDown struct{}
func (a PaneDown) String() string { return "pane down" }
func (a PaneDown) Apply(e *editor.Editor) error {
	e.Down()
	return nil
}
type PaneLeft struct{}

func (a PaneLeft) String() string { return "pane left" }
func (a PaneLeft) Apply(e *editor.Editor) error {
	e.Left()
	return nil
}
type PaneRight struct{}

func (a PaneRight) String() string { return "pane right" }
func (a PaneRight) Apply(e *editor.Editor) error {
	e.Right()
	return nil
}

// command
type CommandRun struct{}

func (a CommandRun) String() string { return "run" }
func (a CommandRun) Apply(e *editor.Editor) error {
	if e.CW.Locked() {
		return nil
	}
	_, err := e.CW.GetCommand()
	if err != nil {
		e.M.SetMode(mode.Normal)
	}
	return nil
}

type CommandInsert struct{ c keyboard.Key }

func (a CommandInsert) String() string { return fmt.Sprintf("insert: %s", string(a.c)) }
func (a CommandInsert) Apply(e *editor.Editor) error {
	e.CW.AddInput(a.c)
	return nil
}

type CommandClearOutput struct{}

func (a CommandClearOutput) String() string { return "clear output" }
func (a CommandClearOutput) Apply(e *editor.Editor) error {
	e.CW.ClearOutput()
	return nil
}

type CommandClearInput struct{}
func (a CommandClearInput) String() string { return "clear Input" }
func (a CommandClearInput) Apply(e *editor.Editor) error {
	e.CW.ClearInput()
	return nil
}

// buffer editing
type Insert struct{ c rune }

func (a Insert) String() string { return fmt.Sprintf("insert: %s", string(a.c)) }
func (a Insert) Apply(e *editor.Editor) error {
	// TODO
	return nil
}
