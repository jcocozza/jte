package editor

import (
	"fmt"
	"strings"

	"github.com/jcocozza/jte/pkg/actions"
	"github.com/jcocozza/jte/pkg/buffer"
	commandwindow "github.com/jcocozza/jte/pkg/commandWindow"
	"github.com/jcocozza/jte/pkg/state"
)

type ActionFn func(e *Editor)
type CommandFn func(e *Editor, args []string) error

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
	actions.DeleteLine:         func(e *Editor) { e.BM.Current.Buf.DeleteLine() },

	actions.Mode_Insert: func(e *Editor) {
		if e.BM.Current.Buf.ReadOnly {
			e.CW.Push(fmt.Sprintf("error: %s is readonly", e.BM.Current.Buf.Name))
			return
		}
		e.SM.SetMode(state.Insert)
	},
	actions.Mode_Normal:  func(e *Editor) { e.SM.SetMode(state.Normal) },
	actions.Mode_Command: func(e *Editor) { e.SM.SetMode(state.Command); e.CW.Activate() },

	// this is handled in the event loop
	actions.InsertCommandChar: nil,
	actions.DeleteCommandChar: func(e *Editor) { e.CW.RemoveCharFromCommand() },
	actions.ClearCommand:      func(e *Editor) { e.CW.Reset() },
	actions.Submit:            func(e *Editor) { e.SM.SetMode(state.Normal) },

	actions.RunCommand: func(e *Editor) {
		cmd, args, err := e.CW.ParseCommand()
		// this means we have an invalid command
		if err != nil {
			e.CW.Push("error: " + err.Error())
			return
		}
		fn, ok := CommandFnRegistry[cmd]
		if !ok {
			panic(fmt.Sprintf("command %d is not registered in the command function registry", cmd))
		}
		err = fn(e, args)
		if err != nil {
			e.CW.Push("error: " + err.Error())
			return
		}
	},
}

// you will usually want to call e.CW.Reset() at the beginning of each of these
var CommandFnRegistry = map[commandwindow.Command]CommandFn{
	commandwindow.LS: func(e *Editor, args []string) error {
		e.CW.Reset()
		bl := e.BM.ListAll()
		for _, bld := range bl {
			e.CW.Push(bld.String())
		}
		e.CW.ShowAll = true
		return nil
	},
	commandwindow.ECHO: func(e *Editor, args []string) error {
		e.CW.Reset()
		toPrint := strings.Join(args, " ")
		e.CW.Push(toPrint)
		return nil
	},
	commandwindow.EDIT: func(e *Editor, args []string) error {
		e.CW.Reset()
		if len(args) == 0 {
			return fmt.Errorf("invalid path")
		}
		path := args[0]
		buf, err := buffer.ReadFileIntoBuffer(path)
		if err != nil {
			return err
		}
		id := e.BM.Add(buf)
		e.BM.SetCurrent(id)
		return nil
	},
	commandwindow.NextBuf: func(e *Editor, args []string) error {
		e.CW.Reset()
		e.BM.Next()
		return nil
	},
	commandwindow.PrevBuf: func(e *Editor, args []string) error {
		e.CW.Reset()
		e.BM.Previous()
		return nil
	},
}
