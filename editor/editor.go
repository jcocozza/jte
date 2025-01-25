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
	bl       []*buffer.Buffer
	currBuf  *buffer.Buffer
	renderer renderer.Renderer
	keyboard *keyboard.Keyboard
	ml       messages.MessageList

	logger *slog.Logger
}

func NewTextEditor(l *slog.Logger) *Editor {
	kb := keyboard.NewKeyboard(l)
	r := &renderer.TextRenderer{}
	return &Editor{
		renderer: r,
		keyboard: kb,
		logger:   l,
	}
}

func (e *Editor) NewBuf() *buffer.Buffer {
	e.currBuf = nil
	newBuf := buffer.NewEmptyBuffer()
	e.bl = append(e.bl, newBuf)
	e.currBuf = newBuf
	return newBuf
}

func (e *Editor) Open(fname string) error {
	b := e.NewBuf()
	return b.Load(fname)
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
		if e.currBuf.C.Y > 0 {
			e.currBuf.C.Y--
		}
	case keyboard.ARROW_DOWN:
		if e.currBuf.C.Y < len(e.currBuf.Rows)-1 {
			e.currBuf.C.Y++
		}
	case keyboard.ARROW_LEFT:
		if e.currBuf.C.X > 0 {
			e.currBuf.C.X--
		}
	case keyboard.ARROW_RIGHT:
		if e.currBuf.C.Y < len(e.currBuf.Rows) && e.currBuf.C.X < len(*e.currBuf.Rows[e.currBuf.C.Y]) {
			e.currBuf.C.X++
		}
	default:
		return nil
	}
	return nil
}

func (e *Editor) PushMessage(msg messages.Message) {
	e.ml.Push(msg)
	e.renderer.SetMsg(e.currBuf, msg)
}

func (e *Editor) openMessages() {
	rows := make([][]byte, len(e.ml)+1)
	rows[0] = []byte("jte messages")
	for i, m := range e.ml {
		rows[i+1] = []byte(fmt.Sprintf("%s - %s", m.Time.String(), m.Text))
	}
	msgBuf := e.NewBuf()
	msgBuf.LoadFromBytes(rows)
}

func (e *Editor) Run() {
	err := e.renderer.Init(e.logger)
	if err != nil {
		panic(err)
	}
	defer e.renderer.Cleanup()
	e.PushMessage(momentoMori)
	for {
		e.renderer.Render(e.currBuf)
		err := e.processKey()
		if err != nil {
			break
		}
	}
}
