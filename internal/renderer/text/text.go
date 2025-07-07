package text

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/jcocozza/jte/internal/buffer"
	"github.com/jcocozza/jte/internal/editor"
	"github.com/jcocozza/jte/internal/panemanager"
	"github.com/jcocozza/jte/internal/renderer/term"
)

type LayoutRect struct {
	X, Y       int
	Rows, Cols int
}

type TextRenderer struct {
	rw *term.RawMode

	screenrows int
	//initscreenrows int
	screencols int

	rowoffset int
	coloffset int

	br *TextBufferRenderer

	// the content that is actually rendered to the screen
	abuf abuf

	logger *slog.Logger
}

func NewTextRenderer(l *slog.Logger) *TextRenderer {
	return &TextRenderer{
		abuf:   abuf{},
		br:     NewTextPaneRenderer(l),
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
	r.drawCursor(offsetY+y - r.br.rowoffset, offsetX+actualCol+1 - r.br.coloffset)
}

func (r *TextRenderer) RenderPane(pn *panemanager.PaneNode, rect LayoutRect, screen [][]byte) {
	if pn == nil {
		return
	}
	switch pn.Direction {
	case panemanager.Horizontal:
		firstH := int(float64(rect.Rows) * pn.Ratio)
		secondH := rect.Rows - firstH
		r.RenderPane(pn.First, LayoutRect{rect.X, rect.Y, firstH, rect.Cols}, screen)
		r.RenderPane(pn.Second, LayoutRect{rect.X, rect.Y + firstH, secondH, rect.Cols}, screen)
		return
	case panemanager.Vertical:
		firstW := int(float64(rect.Cols) * pn.Ratio)
		secondW := rect.Cols - firstW
		r.RenderPane(pn.First, LayoutRect{rect.X, rect.Y, rect.Rows, firstW}, screen)

		for i := 0; i < rect.Rows; i++ {
			screen[i][rect.X+firstW] = '|'
			screen[i][rect.X+firstW+1] = ' '
		}

		r.RenderPane(pn.Second, LayoutRect{rect.X + firstW + 2, rect.Y, rect.Rows, secondW}, screen)
		return
	case panemanager.None:
		rendered := r.br.render(rect.Rows, rect.Cols, pn.Bn.Buf)
		for i := 0; i < len(rendered) && i+rect.Y < len(screen); i++ {
			copy(screen[i+rect.Y][rect.X:], rendered[i])
		}
		return
	default:
		panic("invalid render pane state")
	}
}

func (r *TextRenderer) RenderLayout(root *panemanager.PaneNode, screenrows int, screencols int) [][]byte {
	screen := make([][]byte, screenrows)
	for i := range screen {
		screen[i] = make([]byte, screencols)
		for j := range screen[i] {
			screen[i][j] = ' '
		}
	}
	r.RenderPane(root, LayoutRect{X: 0, Y: 0, Rows: screenrows, Cols: screencols}, screen)
	return screen
}

func (r *TextRenderer) Render(e *editor.Editor) {
	r.logger.Debug("begin rendering")
	r.abuf.Append([]byte("\x1b[?25l")) // hide cursor
	r.abuf.Append([]byte("\x1b[2J"))   // clear entire screen
	r.abuf.Append([]byte("\x1b[H"))    // cursor to home

	rows, cols, _ := r.rw.WindowSize()
	content := r.RenderLayout(e.PM.Root, rows, cols)
	for _, row := range content {
		r.abuf.Append(row)
		//r.abuf.Append([]byte("\x1b[K"))
	}
	r.drawCursorOnBuffer(0, 0, e.PM.Curr.Bn.Buf)
	r.abuf.Flush()
	r.logger.Debug("end rendering")
}
