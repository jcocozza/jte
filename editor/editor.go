package editor

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/jcocozza/jte/api/buffer"
	"github.com/jcocozza/jte/api/keyboard"
	"github.com/jcocozza/jte/api/messages"
	"github.com/jcocozza/jte/api/mode"
	"github.com/jcocozza/jte/api/renderer"
)

type Editor struct {
	bm       *buffer.BufferManager
	renderer renderer.Renderer
	keyboard *keyboard.Keyboard
	ml       messages.MessageList
	mm       *mode.ModeManager
	km       *KeypressManager

	logger *slog.Logger
}

func NewTextEditor(l *slog.Logger) *Editor {
	kb := keyboard.NewKeyboard(l)
	r := &renderer.TextRenderer{}
	bm := buffer.NewBufferManager()
	mm := mode.NewModeManager(l)
	km := NewKeypressManager(l)
	return &Editor{
		bm:       bm,
		renderer: r,
		keyboard: kb,
		logger:   l,
		mm:       mm,
		km:       km,
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
	switch e.mm.Mode() {
	case mode.ModeNavigation:
		e.km.ProcessKeyModeNavigation(e, kp)
		return nil
	case mode.ModeInsert:
		e.km.ProcessKeyModeInsert(e, kp)
		return nil
	}
	return nil
}

func (e *Editor) PushMessage(msg messages.Message) {
	e.ml.Push(msg)
	e.renderer.SetMsg(e.statusInfo(), e.bm.CurrBufNode.Buf, msg)
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

func (e *Editor) statusInfo() renderer.StatusInfo {
	return renderer.StatusInfo{
		Name:      e.bm.CurrBufNode.Buf.Name(),
		Dirty:     e.bm.CurrBufNode.Buf.Dirty(),
		CurrRow:   e.bm.CurrBufNode.Buf.Y(),
		TotalRows: e.bm.CurrBufNode.Buf.TotalRows(),
		Mode:      string(e.mm.Mode()),
	}
}

func (e *Editor) Run() {
	err := e.renderer.Init(e.logger)
	if err != nil {
		panic(err)
	}
	defer e.renderer.Cleanup()
	e.PushMessage(messages.MomentoMori)
	for {
		e.renderer.Render(e.bm.CurrBufNode.Buf, e.statusInfo())
		err := e.processKey()
		if err != nil {
			break
		}
	}
}
