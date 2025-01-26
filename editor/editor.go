package editor

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/jcocozza/jte/api/buffer"
	"github.com/jcocozza/jte/api/keyboard"
	"github.com/jcocozza/jte/api/messages"
	"github.com/jcocozza/jte/api/renderer"
)

type Editor struct {
	bm       *buffer.BufferManager
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
		bm:       bm,
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
	newBuf := buffer.NewEmptyBuffer(fname, e.logger)
	//newBuf, err := buffer.NewLazyBuffer(fname, buffer.BufChunkSize, e.logger)
	//if err != nil {
	//	panic(err)
	//}
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
		e.bm.CurrBufNode.Buf.InsertChar(byte(kp))
		return nil
	}
	switch kp {
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
	default: // if we do not handle the special key explicity, do nothing
		return nil
	}
	return nil
}

func (e *Editor) PushMessage(msg messages.Message) {
	e.ml.Push(msg)
	e.renderer.SetMsg(e.bm.CurrBufNode.Buf, msg)
}

func (e *Editor) openMessages() {
	rows := make([][]byte, len(e.ml)+2)
	rows[0] = []byte("jte messages")
	rows[1] = []byte("------------")
	for i, m := range e.ml {
		rows[i+2] = []byte(fmt.Sprintf("%s - %s", m.Time.String(), m.Text))
	}
	now := time.Now()
	nowStr := now.Format("2006-01-02")
	s := fmt.Sprintf("messages-%s", nowStr)
	newBuf := buffer.NewEmptyBuffer(s, e.logger)
	newBuf.LoadFromBytes(rows)
	id := e.bm.Add(newBuf)
	e.bm.SetCurrent(id)
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
