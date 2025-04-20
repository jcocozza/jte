package renderer

import (
	"bytes"
	"fmt"
	"log/slog"
	"os"

	"github.com/jcocozza/jte/pkg/buffer"
	"github.com/jcocozza/jte/pkg/editor"
	"github.com/jcocozza/jte/pkg/term"
)

const TAB_STOP = 8

type Renderer interface {
	Setup() error
	Exit(msg string)
	ExitErr(err error)
	Render(e *editor.Editor)
}

type TextRenderer struct {
	rw *term.RawMode

	screenrows     int
	initscreenrows int
	screencols     int

	rowoffset int
	coloffset int

	// the content that is actually rendered to the screen
	abuf abuf

	logger *slog.Logger
}

func NewTextRenderer(l *slog.Logger) *TextRenderer {
	return &TextRenderer{
		abuf:   abuf{},
		logger: l.WithGroup("renderer"),
	}
}

func (r *TextRenderer) Setup() error {
	rw, err := term.EnableRawMode()
	if err != nil {
		return err
	}
	r.rw = rw
	rows, cols, err := rw.WindowSize()
	if err != nil {
		return err
	}
	r.screenrows = rows - 1     // leave room for status bar
	r.initscreenrows = rows - 1 // leave room for status bar
	r.screencols = cols
	return nil
}

func (r *TextRenderer) cleanup() {
	if r.rw == nil {
		return
	}
	err := r.rw.Restore()
	if err != nil {
		r.logger.Error("failed to restore terminal", "error", err)
	}
}

func (r *TextRenderer) ExitErr(err error) {
	r.cleanup()
	fmt.Fprintln(os.Stderr, err.Error())
	os.Exit(1)
}

func (r *TextRenderer) Exit(msg string) {
	r.cleanup()
	if msg == "" {
		os.Exit(0)
	}
	fmt.Fprintln(os.Stdout, msg)
	os.Exit(0)
}

func (r *TextRenderer) drawCursor(buf buffer.Buffer) {
	y := (buf.Y() - r.rowoffset) + 1
	actualCol := 0
	for i := 0; i < buf.X(); i++ {
		if buf.Rows[buf.Y()][i] == '\t' {
			actualCol += TAB_STOP - (actualCol % TAB_STOP)
		} else {
			actualCol++
		}
	}
	s := fmt.Sprintf("\x1b[%d;%dH", y, actualCol+1)
	r.abuf.Append([]byte(s))
}

func (r *TextRenderer) drawRow(row []byte) {
	var expanded []byte
	col := 0
	for _, b := range row {
		if b == '\t' {
			spaces := TAB_STOP - (col % TAB_STOP)
			expanded = append(expanded, bytes.Repeat([]byte(" "), spaces)...)
			col += spaces
		} else {
			expanded = append(expanded, b)
			col++
		}
	}
	//line := string(expanded)
	r.abuf.Append(expanded)
}

func (r *TextRenderer) drawBuffer(buf buffer.Buffer) {
	for i := 0; i < r.screenrows; i++ {
		filerow := i + r.rowoffset
		if filerow >= len(buf.Rows) {
			r.abuf.Append([]byte("~"))
		} else {
			r.drawRow(buf.Rows[filerow])
		}
		r.abuf.Append([]byte("\x1b[K"))
		r.abuf.Append([]byte("\r\n"))
	}
}

func (r *TextRenderer) scroll(buf buffer.Buffer) {
	if buf.Y() < r.rowoffset {
		r.rowoffset = buf.Y()
	}
	if buf.Y() >= r.rowoffset+r.screenrows {
		r.rowoffset = buf.Y() - r.screenrows + 1
	}
	if buf.X() < r.coloffset {
		r.coloffset = buf.X()
	}
	if buf.X() >= r.coloffset+r.screencols {
		r.coloffset = buf.X() - r.screencols + 1
	}
}

func (r *TextRenderer) construct(e *editor.Editor) {
	r.scroll(e.BM.Current.Buf)

	r.abuf.Append([]byte("\x1b[?25l")) // hide cursor
	r.abuf.Append([]byte("\x1b[2J"))   // clear entire screen
	r.abuf.Append([]byte("\x1b[H"))    // cursor to home

	r.drawBuffer(e.BM.Current.Buf)
	r.drawCursor(e.BM.Current.Buf)
	r.abuf.Append([]byte("\x1b[?25h")) // show the cursor
}

func (r *TextRenderer) Render(e *editor.Editor) {
	r.construct(e)
	r.abuf.Flush()
}
