package action

import (
	"errors"
	"fmt"

	"github.com/jcocozza/jte/internal/editor"
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
	e.PM.Delete()
	return nil
}
