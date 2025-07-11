package editor

import (
	"errors"
	"fmt"

	"github.com/jcocozza/jte/internal/buffer"
	"github.com/jcocozza/jte/internal/mode"
)

// an action is something that is done to the editor
type Action interface {
	// for debugging
	String() string
	Apply(e *Editor) error
}

// exit

var ErrExit = errors.New("exit")

type Exit struct{}

func (a Exit) String() string { return "exit" }
func (a Exit) Apply(e *Editor) error {
	return ErrExit
}

// modality

type SwitchMode struct {
	m mode.Mode
}

func (a SwitchMode) String() string { return fmt.Sprintf("switch mode: %s", a.m) }
func (a SwitchMode) Apply(e *Editor) error {
	switch a.m {
	case mode.Insert:
		e.BM.Current.Buf.StartEvent(buffer.Event_Insert)
	case mode.Normal:
		e.BM.Current.Buf.Commit()
	case mode.Command:

	default:
		panic("nothing to do there")
	}
	e.m.SetMode(a.m)
	return nil
}

// navigation

type CursorUp struct{}

func (a CursorUp) String() string        { return "CursorUp" }
func (a CursorUp) Apply(e *Editor) error { e.BM.Current.Buf.Up(); return nil }

type CursorDown struct{}

func (a CursorDown) String() string        { return "CursorDown" }
func (a CursorDown) Apply(e *Editor) error { e.BM.Current.Buf.Down(); return nil }

type CursorLeft struct{}

func (a CursorLeft) String() string        { return "CursorLeft" }
func (a CursorLeft) Apply(e *Editor) error { e.BM.Current.Buf.Left(); return nil }

type CursorRight struct{}

func (a CursorRight) String() string        { return "CursorRight" }
func (a CursorRight) Apply(e *Editor) error { e.BM.Current.Buf.Right(); return nil }

// splits

type SplitVertical struct{}

func (a SplitVertical) String() string { return "vert split" }
func (a SplitVertical) Apply(e *Editor) error {
	e.Active.Pane.Active = false
	e.Active = e.Active.SplitVertical()
	e.Active.Pane.Active = true
	return nil
}

type SplitHorizontal struct{}

func (a SplitHorizontal) String() string { return "horizontal split" }
func (a SplitHorizontal) Apply(e *Editor) error {
	e.Active.Pane.Active = false
	e.Active = e.Active.SplitHorizontal()
	e.Active.Pane.Active = true
	return nil
}

// buffer stuff

type Commit struct{}

func (a Commit) String() string        { return "commit" }
func (a Commit) Apply(e *Editor) error { e.BM.Current.Buf.Commit(); return nil }

type Insert struct{ c rune }

func (a Insert) String() string { return fmt.Sprintf("insert: %s", string(a.c)) }
func (a Insert) Apply(e *Editor) error {
	c := buffer.Insert{Contents: [][]rune{{a.c}}}
	return e.BM.Current.Buf.AcceptChange(c)
}

type EnterNewLine struct{}

func (a EnterNewLine) String() string { return "new line (enter)" }
func (a EnterNewLine) Apply(e *Editor) error {
	c := buffer.EnterNewLine{}
	return e.BM.Current.Buf.AcceptChange(c)
}

type NewLineAbove struct{}

func (a NewLineAbove) String() string { return "new line (above)" }
func (a NewLineAbove) Apply(e *Editor) error {
	y := e.BM.Current.Buf.Y() - 1
	if y < 0 {
		y = 0
	}
	c := buffer.InsertNewLine{Y: y}
	return e.BM.Current.Buf.AcceptChange(c)
}

type NewLineBelow struct{}

func (a NewLineBelow) String() string { return "new line (below)" }
func (a NewLineBelow) Apply(e *Editor) error {
	y := e.BM.Current.Buf.Y() + 1
	if y > len(e.BM.Current.Buf.Rows) {
		y = len(e.BM.Current.Buf.Rows)
	}
	c := buffer.InsertNewLine{Y: y}
	return e.BM.Current.Buf.AcceptChange(c)
}

type Backspace struct{}

func (a Backspace) String() string { return "backspace" }
func (a Backspace) Apply(e *Editor) error {
	c := buffer.Backspace{}
	return e.BM.Current.Buf.AcceptChange(c)
}

type Delete struct{}

func (a Delete) String() string { return "Delete" }
func (a Delete) Apply(e *Editor) error {
	c := buffer.Delete{}
	return e.BM.Current.Buf.AcceptChange(c)
}

type DeleteLine struct{}

func (a DeleteLine) String() string { return "DeleteLine" }
func (a DeleteLine) Apply(e *Editor) error {
	c := buffer.DeleteLine{}
	return e.BM.Current.Buf.StartAndAcceptChange(c, buffer.Event_Delete)
}

// command stuff
type InsertCommandChar struct{ c rune }

func (a InsertCommandChar) String() string        { return fmt.Sprintf("insert command char %s", string(a.c)) }
func (a InsertCommandChar) Apply(e *Editor) error { return nil }
