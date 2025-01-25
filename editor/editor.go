package editor

import (
	"fmt"
	"log/slog"

	"github.com/jcocozza/jte/api/buffer"
	"github.com/jcocozza/jte/api/keyboard"
	"github.com/jcocozza/jte/api/messages"
	"github.com/jcocozza/jte/api/renderer"
)

type Editor struct {
	bm *buffer.BufferManager
	renderer renderer.Renderer
	keyboard *keyboard.Keyboard
	ml       messages.MessageList

	logger *slog.Logger
}

func NewTextEditor(l *slog.Logger) *Editor {
	kb := keyboard.NewKeyboard(l)
	r := &renderer.TextRenderer{}
	bm := buffer.NewBufferManager()
	return &Editor{
		bm: bm,
		renderer: r,
		keyboard: kb,
		logger:   l,
	}
}

// create a new buffer and set it to the current
func (e *Editor) NewBuf() buffer.Buffer {
	newBuf := buffer.NewEmptyBuffer("", e.logger)
	id := e.bm.Add(newBuf)
	e.bm.SetCurrent(id)
	return newBuf
}

func (e *Editor) Open(fname string) error {
	//b := e.NewBuf()
	//newBuf := buffer.NewEmptyBuffer(fname, e.logger)
	newBuf, err := buffer.NewLazyBuffer(fname, buffer.BufChunkSize, e.logger)
	if err != nil {
		panic(err)
	}
	id := e.bm.Add(newBuf)
	e.bm.SetCurrent(id)
	return newBuf.Load()
}

func (e *Editor) processKey() error {
	kp, err := e.keyboard.GetKeypress()
	if err != nil {
		return err
	}
	if kp.IsUnicode() {
		return nil
	}
	switch kp.Key {
	case keyboard.CtrlQ:
		e.renderer.Exit("regular quit")
	case keyboard.CtrlL:
		e.openMessages()
	case keyboard.ARROW_UP:
		e.bm.CurrBufNode.Buf.Up()
	case keyboard.ARROW_DOWN:
		e.bm.CurrBufNode.Buf.Down()
	case keyboard.ARROW_LEFT:
		e.bm.CurrBufNode.Buf.Left()
	case keyboard.ARROW_RIGHT:
		e.bm.CurrBufNode.Buf.Right()
	default:
		return nil
	}
	return nil
}

func (e *Editor) PushMessage(msg messages.Message) {
	e.ml.Push(msg)
	e.renderer.SetMsg(e.bm.CurrBufNode.Buf, msg)
}

func (e *Editor) openMessages() {
	rows := make([][]byte, len(e.ml)+1)
	rows[0] = []byte("jte messages")
	for i, m := range e.ml {
		rows[i+1] = []byte(fmt.Sprintf("%s - %s", m.Time.String(), m.Text))
	}
	// TODO: re-implement
	//msgBuf := e.NewBuf()
	//msgBuf.LoadFromBytes(rows)
}

func (e *Editor) Run() {
	err := e.renderer.Init(e.logger)
	if err != nil {
		panic(err)
	}
	defer e.renderer.Cleanup()
	e.PushMessage(messages.MomentoMori)
	for {
		e.renderer.Render(e.bm.CurrBufNode.Buf)
		err := e.processKey()
		if err != nil {
			break
		}
	}
}
