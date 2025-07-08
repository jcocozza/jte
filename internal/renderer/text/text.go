package text

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/jcocozza/jte/internal/buffer"
	"github.com/jcocozza/jte/internal/editor"
	"github.com/jcocozza/jte/internal/mode"
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

	currRowoffset int
	currColoffset int
	currRect      LayoutRect

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
	s := fmt.Sprintf("\x1b[%d;%dH", y+1, x+1) // set cursor position
	r.abuf.Append([]byte(s))
	r.abuf.Append([]byte("\x1b[?25h")) // show cursor
}

func (r *TextRenderer) renderCursor(x int, y int, currRow buffer.BufRow) {
	// unsure if i need this - i think i might later
	// for now,  just leaving it commented out
	//actualCol := 0
	//for _, b := range currRow {
	//	if b == '\t' {
	//		actualCol += TAB_STOP - (actualCol % TAB_STOP)
	//	} else {
	//		actualCol++
	//	}
	//}

	rx := x - r.currColoffset + r.currRect.X
	ry := y - r.currRowoffset + r.currRect.Y
	maxCursorRow := r.currRect.Y + r.currRect.Rows - 2
	r.logger.Debug("render cursor",
		slog.Int("rendered x", rx),
		slog.Int("real x", x),
		slog.Int("rendered y", ry),
		slog.Int("real y", y),
		slog.Int("col offset", r.currColoffset),
		slog.Int("row offset", r.currRowoffset),
		slog.Int("max row", maxCursorRow),
	)
	//if ry > maxCursorRow {
	//	ry  = maxCursorRow
	//}
	//if ry < r.currRect.Y {
	//	ry = r.currRect.Y
	//} else if ry > maxCursorRow {
	//	ry = maxCursorRow
	//}
	if ry > maxCursorRow {
		ry = maxCursorRow
	}

	r.drawCursor(rx, ry)
}

// TODO: need a way to determine if the panenode is the current
func (r *TextRenderer) RenderPane(pn *panemanager.PaneNode, es *editor.EditorStatus, rect LayoutRect, screen [][]byte) {
	if pn == nil {
		return
	}
	switch pn.Direction {
	case panemanager.Horizontal:
		firstH := int(float64(rect.Rows) * pn.Ratio)
		secondH := rect.Rows - firstH
		r.RenderPane(pn.First, es, LayoutRect{rect.X, rect.Y, firstH, rect.Cols}, screen)
		r.RenderPane(pn.Second, es, LayoutRect{rect.X, rect.Y + firstH, secondH, rect.Cols}, screen)
		return
	case panemanager.Vertical:
		firstW := int(float64(rect.Cols) * pn.Ratio)
		secondW := rect.Cols - firstW
		r.RenderPane(pn.First, es, LayoutRect{rect.X, rect.Y, rect.Rows, firstW}, screen)

		for i := 0; i < rect.Rows-1; i++ {
			screen[i][rect.X+firstW] = '|'
			screen[i][rect.X+firstW+1] = ' '
		}

		r.RenderPane(pn.Second, es, LayoutRect{rect.X + firstW + 2, rect.Y, rect.Rows, secondW}, screen)
		return
	case panemanager.None:
		r.logger.Debug("rect", slog.Any("rect", rect))
		rendered := r.br.render(rect.Rows, rect.Cols, pn.Bn.Buf)
		for i := 0; i < len(rendered) && i+rect.Y < len(screen); i++ {
			copy(screen[i+rect.Y][rect.X:], rendered[i])
		}
		/*
			THE STATUS BAR IS DRAWN HERE
		*/
		var active string
		if es.CurrentPane == pn {
			active = "(a)"
			r.currRect = rect
			r.currColoffset = r.br.coloffset
			r.currRowoffset = r.br.rowoffset
		} else {
			active = ""
		}

		status := fmt.Appendf([]byte{}, "%s[%s] %s (%s) ln: %d/%d", active, es.Mode.String(), pn.Bn.Buf.Name, pn.Bn.Buf.FileType.String(), pn.Bn.Buf.Y(), len(pn.Bn.Buf.Rows))
		copy(screen[len(rendered)-1+rect.Y][rect.X:], status)
		return
	default:
		panic("invalid render pane state")
	}
}

func (r *TextRenderer) RenderLayout(root *panemanager.PaneNode, es *editor.EditorStatus, screenrows int, screencols int) [][]byte {
	screen := make([][]byte, screenrows)
	for i := range screen {
		screen[i] = make([]byte, screencols)
		for j := range screen[i] {
			screen[i][j] = ' '
		}
	}
	r.RenderPane(root, es, LayoutRect{X: 0, Y: 0, Rows: screenrows, Cols: screencols}, screen)
	return screen
}

func (r *TextRenderer) Render(e *editor.Editor) {
	r.logger.Debug("begin rendering")
	r.abuf.Append([]byte("\x1b[?25l")) // hide cursor
	r.abuf.Append([]byte("\x1b[2J"))   // clear entire screen
	r.abuf.Append([]byte("\x1b[H"))    // cursor to home

	rows, cols, _ := r.rw.WindowSize()
	//rows -= len(e.CW.Output) //+ 1 // leave room for command line and possible output

	if len(e.CW.Output) > 0 {
		rows -= len(e.CW.Output)
	} else {
		rows -= 1
	}

	content := r.RenderLayout(e.PM.Root, e.Status(), rows, cols)
	for _, row := range content {
		r.abuf.Append(row)
		//r.abuf.Append([]byte("\x1b[K"))
	}

	// render command window output
	if e.CW.ShowOutput {
		for i, o := range e.CW.Output {
			r.abuf.Append([]byte(o))
			if i < len(e.CW.Output)-1 {
				r.abuf.Append([]byte("\n"))
			}
		}
	} else if e.M.Current() == mode.Command {
		r.abuf.Append([]byte("> "))
		r.abuf.Append([]byte(string(e.CW.Input)))
	}

	x, y := e.BM.Current.Buf.X(), e.BM.Current.Buf.Y()
	r.logger.Debug("curr rect", slog.Any("rect", r.currRect))

	//r.logger.Debug("current buf offsets", slog.Int("col offset", r.coloffset), slog.Int("row offset", r.rowoffset))
	r.renderCursor(x+r.br.gutterShift, y, e.BM.Current.Buf.Rows[e.BM.Current.Buf.Y()])
	r.abuf.Flush()
	r.logger.Debug("end rendering")
}
