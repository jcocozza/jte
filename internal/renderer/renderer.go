package renderer

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/jcocozza/jte/internal/buffer"
	"github.com/jcocozza/jte/internal/editor"
	"github.com/jcocozza/jte/internal/term"
)

type Renderer interface {
	Setup() error
	Exit(msg string)
	ExitErr(err error)
	Render(e *editor.Editor)
}

type TextRenderer struct {
	rw *term.RawMode

	screenrows int
	//initscreenrows int
	screencols int

	rowoffset int
	coloffset int

	lr *LayoutRenderer
	pr *TextPaneRenderer

	// the content that is actually rendered to the screen
	abuf abuf

	logger *slog.Logger
}

func NewTextRenderer(l *slog.Logger) *TextRenderer {
	return &TextRenderer{
		abuf:   abuf{},
		lr:     NewLayoutRenderer(l),
		pr:     NewTextPaneRenderer(l),
		logger: l.WithGroup("renderer"),
	}
}

func (r *TextRenderer) Setup() error {
	rw, err := term.EnableRawMode()
	if err != nil {
		return err
	}
	r.rw = rw
	return nil
}

func (r *TextRenderer) cleanup() {
	r.abuf.Append([]byte("\x1b[2J")) // clear entire screen
	r.abuf.Flush()
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

func (r *TextRenderer) drawCursor(x int, y int) {
	s := fmt.Sprintf("\x1b[%d;%dH", x, y) // set cursor position
	r.abuf.Append([]byte(s))
	r.abuf.Append([]byte("\x1b[?25h")) // show cursor
}

func (r *TextRenderer) drawCursorOnBuffer(offsetX int, offsetY int, buf *buffer.Buffer) {
	y := (buf.Y() - r.rowoffset) + 1
	actualCol := 0
	for i := 0; i < buf.X(); i++ {
		if buf.Rows[buf.Y()][i] == '\t' {
			actualCol += TAB_STOP - (actualCol % TAB_STOP)
		} else {
			actualCol++
		}
	}
	r.drawCursor(offsetY+y, offsetX+actualCol+1)
}

func (r *TextRenderer) Render(e *editor.Editor) {
	r.logger.Debug("begin rendering")
	r.abuf.Append([]byte("\x1b[?25l")) // hide cursor
	r.abuf.Append([]byte("\x1b[2J"))   // clear entire screen
	r.abuf.Append([]byte("\x1b[H"))    // cursor to home

	rows, cols, _ := r.rw.WindowSize()
	content := r.lr.RenderLayout(e.Root, r.pr, rows, cols)
	for _, row := range content {
		r.logger.Log(context.TODO(), slog.LevelDebug-1, "row", slog.String("row", string(row)))
		r.abuf.Append(row)
		//r.abuf.Append([]byte("\x1b[K"))
	}

	r.drawCursorOnBuffer(0,0, e.Active.Pane.Buf)
	r.abuf.Flush()
	r.logger.Debug("end rendering")
}
