package editor

import (
	"github.com/jcocozza/jte/pkg/actions"
)

type ActionFn func(e *Editor)

var Registry = map[actions.Action]ActionFn{
	actions.None:        func(e *Editor) {},
	actions.Exit:        func(e *Editor) {},
	actions.CursorUp:    func(e *Editor) { e.BM.Current.Buf.Up() },
	actions.CursorDown:  func(e *Editor) { e.BM.Current.Buf.Down() },
	actions.CursorLeft:  func(e *Editor) { e.BM.Current.Buf.Left() },
	actions.CursorRight: func(e *Editor) { e.BM.Current.Buf.Right() },
}
