package editor

import (
	"github.com/jcocozza/jte/pkg/actions"
	"github.com/jcocozza/jte/pkg/state"
)

type ActionFn func(e *Editor)

var Registry = map[actions.Action]ActionFn{
	actions.None:   func(e *Editor) {},
	actions.Exit:   func(e *Editor) {},
	actions.Repeat: func(e *Editor) {},

	actions.CursorUp:    func(e *Editor) { e.BM.Current.Buf.Up() },
	actions.CursorDown:  func(e *Editor) { e.BM.Current.Buf.Down() },
	actions.CursorLeft:  func(e *Editor) { e.BM.Current.Buf.Left() },
	actions.CursorRight: func(e *Editor) { e.BM.Current.Buf.Right() },
	actions.StartLine:   func(e *Editor) { e.BM.Current.Buf.StartLine() },
	actions.EndLine:     func(e *Editor) { e.BM.Current.Buf.EndLine() },

	// this is handled in the event loop
	actions.InsertChar:         nil,
	actions.InsertNewLine:      func(e *Editor) { e.BM.Current.Buf.InsertNewLine() },
	actions.InsertNewLineAbove: func(e *Editor) { e.BM.Current.Buf.InsertNewLineAbove() },
	actions.InsertNewLineBelow: func(e *Editor) { e.BM.Current.Buf.InsertNewLineBelow() },
	actions.DeleteChar:         func(e *Editor) { e.BM.Current.Buf.DeleteChar() },
	actions.RemoveChar:         func(e *Editor) { panic("remove char; unimplemented") },
	actions.DeleteLine:         func(e *Editor) { panic("delete line: unimplemented") },

	actions.Mode_Insert: func(e *Editor) { e.SM.SetMode(state.Insert) },
	actions.Mode_Normal: func(e *Editor) { e.SM.SetMode(state.Normal) },
}
