package editor

import (
	"log/slog"

	"github.com/jcocozza/jte/api/command"
	"github.com/jcocozza/jte/api/keyboard"
	"github.com/jcocozza/jte/api/mode"
	//"github.com/jcocozza/jte/api/renderer"
)

/*
for now this is in the editor package

i may mode it elsewhere if that makes sense, but currently not sure how i want to handle modality correctly
*/

type KeypressBuf []keyboard.Key

func (k *KeypressBuf) Append(key keyboard.Key) {
	*k = append(*k, key)
}

func (k *KeypressBuf) Clear() {
	*k = []keyboard.Key{}
}

// this is the naive approach to a modal editor
// in the future, i plan to improve this
type KeypressManager struct {
	kb KeypressBuf

	logger *slog.Logger
}

func NewKeypressManager(l *slog.Logger) *KeypressManager {
	return &KeypressManager{
		kb:     []keyboard.Key{},
		logger: l,
	}
}

func (k *KeypressManager) ProcessKeyModeNavigation(e *Editor, key keyboard.Key) {
	switch key {
	case keyboard.CtrlQ:
		e.renderer.Exit("regular quit")
	case keyboard.CtrlL:
		e.openMessages()
	case keyboard.ARROW_UP, 'k':
		e.bm.CurrBufNode.Buf.Up()
	case keyboard.ARROW_DOWN, 'j':
		e.bm.CurrBufNode.Buf.Down()
	case keyboard.ARROW_LEFT, 'h':
		e.bm.CurrBufNode.Buf.Left()
	case keyboard.ARROW_RIGHT, 'l':
		e.bm.CurrBufNode.Buf.Right()
	case 'i':
		e.mm.SetMode(mode.ModeInsert)
	case ':':
		e.mm.SetMode(mode.ModeCommand)
		e.cw.Activate()
	case '/', keyboard.CtrlF:
		e.mm.SetMode(mode.ModeCommand)
		e.cw.ActivateSearch()
	default:
		k.kb.Append(key)
	}
}

func (k *KeypressManager) ProcessKeyModeInsert(e *Editor, key keyboard.Key) {
	if key.IsUnicode() {
		e.bm.CurrBufNode.Buf.InsertChar(byte(key))
		return
	}
	switch key {
	case keyboard.CtrlQ:
		e.renderer.Exit("regular quit")
	case keyboard.ARROW_UP:
		e.bm.CurrBufNode.Buf.Up()
	case keyboard.ARROW_DOWN:
		e.bm.CurrBufNode.Buf.Down()
	case keyboard.ARROW_LEFT:
		e.bm.CurrBufNode.Buf.Left()
	case keyboard.ARROW_RIGHT:
		e.bm.CurrBufNode.Buf.Right()
	case keyboard.BACKSPACE, keyboard.BACKSPACE_2:
		e.bm.CurrBufNode.Buf.DeleteChar()
	case keyboard.DELETE:
		// TODO: this logic needs to be better encapsulated in the buffer
		if e.bm.CurrBufNode.Buf.Y() < e.bm.CurrBufNode.Buf.NumRows() && e.bm.CurrBufNode.Buf.X() < len(e.bm.CurrBufNode.Buf.Row(e.bm.CurrBufNode.Buf.Y())) {
			e.bm.CurrBufNode.Buf.Right()
		}
		e.bm.CurrBufNode.Buf.DeleteChar()
	case keyboard.ENTER:
		e.bm.CurrBufNode.Buf.InsertNewLine()
	case keyboard.HOME:
		e.bm.CurrBufNode.Buf.StartLine()
	case keyboard.END:
		e.bm.CurrBufNode.Buf.EndLine()
	case keyboard.TAB:
		e.bm.CurrBufNode.Buf.InsertChar('\t')
		//for i := 0; i < renderer.TAB_STOP; i++ {
		//	e.bm.CurrBufNode.Buf.InsertChar(' ')
		//}
	case keyboard.ESC:
		e.mm.SetMode(mode.ModeNavigation)
	}
}

func (k *KeypressManager) ProcessKeyModeCommand(e *Editor, key keyboard.Key) {
	switch e.cw.Mode {
	case command.CommandInactive:
		return // this shouldn't happen
	case command.CommandBasic:
		switch {
		case key.IsUnicode():
			e.cw.AddInput(key)
			return
		case key == keyboard.BACKSPACE_2:
			e.cw.ShrinkInput()
		case key == keyboard.ESC:
			e.cw.Clear()
			e.mm.SetMode(mode.ModeNavigation)
			e.cw.Mode = command.CommandInactive
			return
		case key == keyboard.ENTER:
			e.cw.Handle(e.bm)
			e.renderer.Render(e.bm.CurrBufNode.Buf, e.statusInfo(), e.cw)
			_, _ = e.keyboard.GetKeypress()
			e.cw.Clear()
			e.renderer.Render(e.bm.CurrBufNode.Buf, e.statusInfo(), e.cw)
			e.mm.SetMode(mode.ModeNavigation)
			e.cw.Mode = command.CommandInactive
			return
		}
	case command.CommandSearch:
		switch {
		case key == keyboard.ARROW_RIGHT:
			loc := e.cw.SearchResults.Next()
			e.bm.CurrBufNode.Buf.GoTo(loc.X, loc.Y)
		case key == keyboard.ARROW_LEFT:
			loc := e.cw.SearchResults.Previous()
			e.bm.CurrBufNode.Buf.GoTo(loc.X, loc.Y)
		case key == keyboard.BACKSPACE_2:
			e.cw.ShrinkInput()
		case key.IsUnicode():
			e.cw.AddInput(key)
			e.cw.HandleSearch(e.bm.CurrBufNode.Buf)
			closest := e.cw.SearchResults.Current()
			e.bm.CurrBufNode.Buf.GoTo(closest.X, closest.Y)
			return
		case key == keyboard.ENTER:
			e.cw.Mode = command.CommandInactive
			e.cw.Clear()
			e.mm.SetMode(mode.ModeNavigation)
			return
		case key == keyboard.ESC:
			e.cw.Mode = command.CommandInactive
			e.cw.Clear()
			e.mm.SetMode(mode.ModeNavigation)
			return
		}
	}
}
