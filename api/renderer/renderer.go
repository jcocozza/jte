package renderer

import (
	"bytes"
	"fmt"
	"log/slog"
	"os"

	"github.com/jcocozza/jte/api/buffer"
	"github.com/jcocozza/jte/api/messages"
	"github.com/jcocozza/jte/term"
)

type Renderer interface {
	Init(l *slog.Logger) error
	Render(buf *buffer.Buffer)
	Cleanup()
	Exit(msg string)
	ExitErr(err error)
	SetMsg(buf *buffer.Buffer, msg messages.Message)
}

type erow struct {
	render []byte
	hl     []int
}

type TextRenderer struct {
	screenrows int
	screencols int

	rowoffset int
	coloffset int

	rw *term.RawMode

	abuf abuf

	currMsg messages.Message

	logger *slog.Logger
}

func (r *TextRenderer) Init(l *slog.Logger) error {
	rw, err := term.EnableRawMode()
	if err != nil {
		return err
	}
	r.rw = rw
	rows, cols, err := rw.WindowSize()
	if err != nil {
		return err
	}
	r.abuf = abuf{}
	r.screenrows = rows - 2 // leave room for status bar and messages
	r.screencols = cols
	r.logger = l
	return nil
}

func (r *TextRenderer) Cleanup() {
	if r.rw != nil {
		r.logger.Info("running cleanup")
		if err := r.rw.Restore(); err != nil {
			r.logger.Error("failed to restore terminal", "error", err)
		}
	}
}

func (r *TextRenderer) ExitErr(err error) {
	r.Cleanup()
	fmt.Fprintln(os.Stderr, err.Error())
	os.Exit(1)
}

func (r *TextRenderer) Exit(msg string) {
	r.Cleanup()
	fmt.Fprintln(os.Stdout, msg)
	os.Exit(0)
}

func (r *TextRenderer) drawCursor(buf *buffer.Buffer) {
	y := (buf.C.Y - r.rowoffset) + 1
	x := (buf.C.X - r.coloffset) + 1
	s := fmt.Sprintf("\x1b[%d;%dH", y, x)
	r.abuf.Append([]byte(s))
}

func (r *TextRenderer) SetMsg(buf *buffer.Buffer, msg messages.Message) {
	r.currMsg = msg
	r.Render(buf)
}

var welcome []byte = []byte("jte -- version: v0.0.1")

func (r *TextRenderer) drawWelcome() {
	for i := 0; i < r.screenrows; i++ {
		if i == r.screenrows/3 {
			r.abuf.Append([]byte("~"))
			padding := (r.screencols - len(welcome)) / 2
			r.abuf.Append(bytes.Repeat([]byte(" "), padding))
			r.abuf.Append(welcome)
		} else {
			r.abuf.Append([]byte("~"))
		}
		r.abuf.Append([]byte("\x1b[K"))
		r.abuf.Append([]byte("\r\n"))
	}
}

func (r *TextRenderer) status(buf *buffer.Buffer) string {
	displayName := buf.Name
	if displayName == "" {
		displayName = "[No Name]"
	}
	var displayDirty string = ""
	if buf.Dirty {
		displayDirty = "(modified)"
	}
	var displayRowNum int = 0
	if len(buf.Rows) != 0 {
		displayRowNum = len(buf.Rows) - 1 // -1 because i want a 0 indexed system
	}
	return fmt.Sprintf("ln:%d/%d - %s %s", buf.C.Y, displayRowNum, displayDirty, displayName)
}

func (r *TextRenderer) drawStatusBar(buf *buffer.Buffer) {
	// status bar
	r.abuf.Append([]byte("\x1b[7m"))
	status := r.status(buf)

	r.abuf.Append(bytes.Repeat([]byte(" "), r.screencols-len(status)))
	r.abuf.Append([]byte(status))
	r.abuf.Append([]byte("\x1b[m"))
	r.abuf.Append([]byte("\r\n"))
	// messages
	r.abuf.Append([]byte("\x1b[K"))
	var displayMsg string = r.currMsg.Text
	if len(r.currMsg.Text) > r.screencols {
		displayMsg = r.currMsg.Text[0:r.screencols]
	}
	if r.currMsg.NonEmpty() && !r.currMsg.Expired() {
		r.abuf.Append([]byte(displayMsg))
	}
}

func (r *TextRenderer) drawBuffer(buf *buffer.Buffer) {
	for i := range r.screenrows {
		filerow := i + r.rowoffset
		if filerow < 0 || filerow >= len(buf.Rows) {
			r.abuf.Append([]byte("~"))
		} else {
			r.logger.Debug("rendering", slog.String("row",string(*buf.Rows[filerow])))
			r.abuf.Append(*buf.Rows[filerow])
		}
		r.abuf.Append([]byte("\x1b[K"))
		r.abuf.Append([]byte("\r\n"))
	}
}

/*
func (r *TextRenderer) Clear() {
	r.abuf.Clear()
	r.abuf.Flush()
	r.rowoffset = 0
	r.coloffset = 0
}
*/

func (r *TextRenderer) scroll(buf *buffer.Buffer) {
	if buf.C.Y < r.rowoffset {
		r.rowoffset = buf.C.Y
	}
	if buf.C.Y >= r.rowoffset+r.screenrows {
		r.rowoffset = buf.C.Y - r.screenrows + 1
	}
	if buf.C.X < r.coloffset {
		r.coloffset = buf.C.X
	}
	if buf.C.X >= r.coloffset+r.screencols {
		r.coloffset = buf.C.X - r.screencols + 1
	}

}

func (r *TextRenderer) Render(buf *buffer.Buffer) {
	r.scroll(buf)
	r.abuf.Append([]byte("\x1b[?25l"))
	r.abuf.Append([]byte("\x1b[H"))

	if buf.IsEmpty() {
		r.drawWelcome()
	} else {
		r.drawBuffer(buf)
	}
	r.drawStatusBar(buf)

	r.drawCursor(buf)
	r.abuf.Append([]byte("\x1b[?25h"))
	r.abuf.Flush()
}
