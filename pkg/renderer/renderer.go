package renderer

import (
	"bytes"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/jcocozza/jte/pkg/buffer"
	"github.com/jcocozza/jte/pkg/editor"
	"github.com/jcocozza/jte/pkg/filetypes"
	"github.com/jcocozza/jte/pkg/state"
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

	screenrows int
	//initscreenrows int
	screencols int

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

func (r *TextRenderer) drawCursorOnBuffer(buf *buffer.Buffer, gutterMaxSize int) {
	y := (buf.Y() - r.rowoffset) + 1
	actualCol := 0
	for i := 0; i < buf.X(); i++ {
		if buf.Rows[buf.Y()][i] == '\t' {
			actualCol += TAB_STOP - (actualCol % TAB_STOP)
		} else {
			actualCol++
		}
	}
	r.drawCursor(y, actualCol+1+gutterMaxSize)
}

func (r *TextRenderer) drawCursor(x int, y int) {
	s := fmt.Sprintf("\x1b[%d;%dH", x, y)
	r.abuf.Append([]byte(s))
}

type RowData struct {
	ftype         filetypes.FileType
	row           []byte
	rownum        int
	gutterMaxSize int
}

func (r *TextRenderer) drawRow(buf *buffer.Buffer, rd RowData) {
	// if i ever decide to do a 'background' color
	//r.abuf.Append([]byte("\x1b[48;5;25m\x1b[38;5;231m"))
	var expanded []byte
	col := 0
	for _, b := range rd.row {
		if b == '\t' {
			spaces := TAB_STOP - (col % TAB_STOP)
			expanded = append(expanded, bytes.Repeat([]byte(" "), spaces)...)
			col += spaces
		} else {
			expanded = append(expanded, b)
			col++
		}
	}
	line := string(expanded)
	tokens := filetypes.GetMatches(rd.ftype, line)

	// quick helper func to compute absolute values
	abs := func(i int) int {
		if i < 0 {
			return i * -1
		}
		return i
	}

	// compute relative line numbers
	relLine := abs(rd.rownum - buf.Y())

	var finalGutterMsg string
	if relLine == 0 {
		gutterMsg := strconv.Itoa(rd.rownum)
		spacing := strings.Repeat(" ", rd.gutterMaxSize-len(gutterMsg))
		finalGutterMsg = gutterMsg + spacing + " "
	} else {
		gutterMsg := strconv.Itoa(abs(rd.rownum - buf.Y()))
		spacing := strings.Repeat(" ", rd.gutterMaxSize-len(gutterMsg))
		finalGutterMsg = spacing + gutterMsg + " "
	}
	r.abuf.Append([]byte(finalGutterMsg))
	//rowLen := 0 // for the background color
	for _, tkn := range tokens {
		r.abuf.Append([]byte(tkn.Color))
		r.abuf.Append([]byte(tkn.Text))
		r.abuf.Append([]byte(filetypes.Reset))
		// rowLen++
	}

	/* if i ever decide to do a 'background' color
	if rowLen < r.screencols {
		_ = strings.Repeat(" ", r.screencols-rowLen)
		r.abuf.Append([]byte("\x1b[48;5;25m\x1b[38;5;231m")) // ensure BG color matches
	}
	r.abuf.Append(expanded)
	*/
}

func (r *TextRenderer) drawBuffer(buf *buffer.Buffer, gutterMaxSize int) {
	for i := 0; i < r.screenrows; i++ {
		filerow := i + r.rowoffset
		if filerow >= len(buf.Rows) {
			r.abuf.Append([]byte("~"))
		} else {
			rd := RowData{
				ftype:         buf.FileType,
				row:           buf.Rows[filerow],
				rownum:        filerow,
				gutterMaxSize: gutterMaxSize,
			}
			r.drawRow(buf, rd)
		}
		r.abuf.Append([]byte("\x1b[K"))
		r.abuf.Append([]byte("\r\n"))
	}
}

func (r *TextRenderer) drawStatusBar(e *editor.Editor) {
	r.abuf.Append([]byte("\x1b[7m"))
	mode := e.SM.Current()

	name := e.BM.Current.Buf.Name

	var displayModified string = ""
	if e.BM.Current.Buf.Modified {
		displayModified = "(Δ)"
	}

	var displayRowNum int = 0
	totalRows := len(e.BM.Current.Buf.Rows)
	currRow := e.BM.Current.Buf.Y()
	if totalRows != 0 {
		displayRowNum = totalRows - 1 // -1 because i want a 0 indexed system
	}
	status := fmt.Sprintf("ln:%d/%d - %s %s", currRow, displayRowNum, displayModified, name)

	r.abuf.Append([]byte(mode))
	r.abuf.Append(bytes.Repeat([]byte(" "), r.screencols-len(status)-len(mode)))
	r.abuf.Append([]byte(status))
	r.abuf.Append([]byte("\x1b[m"))
	//r.abuf.Append([]byte("\r\n"))
	r.abuf.Append([]byte("\x1b[K"))
}

func (r *TextRenderer) setupThisRender(e *editor.Editor) {
	rows, cols, _ := r.rw.WindowSize()

	//largestRowNumLen := len(strconv.Itoa(cols))

	r.screenrows = rows - 1 // leave room for status bar and command message
	r.screencols = cols     //- largestRowNumLen // leave room for line nums

	// add the 0 check to ensure we always have at least 1 row for the message
	if e.CW.ShowAll && e.CW.Size() != 0 {
		r.screenrows = r.screenrows - e.CW.Size()
	} else {
		r.screenrows-- // otherwise just leave room for 1 message
	}
}

func (r *TextRenderer) drawCommandArea(e *editor.Editor) {
	// we don't care if there are no messages
	// add the 0 check to ensure we always have at least 1 row for the message
	if e.CW.ShowAll && e.CW.Size() != 0 {
		contents := e.CW.Dump()
		last := len(contents) - 1
		for i, c := range contents {
			r.abuf.Append([]byte(c))
			if i != last {
				r.abuf.Append([]byte("\r\n"))
			}
		}
		return
	}
	// if the current command is not empty,
	// then the user is typing
	// so we want to keep rendering that
	if e.CW.Active() {
		r.abuf.Append([]byte("> "))
		r.abuf.Append([]byte(e.CW.CmdBuf()))
		return
	}
	// otherwise just show the next message
	// (if there is one)
	r.abuf.Append([]byte("\r\n"))
	msg := e.CW.Next()
	r.abuf.Append([]byte(msg))
}

func (r *TextRenderer) scroll(buf *buffer.Buffer) {
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
	gutterMaxSize := len(strconv.Itoa(len(e.BM.Current.Buf.Rows))) + 1 // +1 for a space after the gutter content
	r.setupThisRender(e)
	r.scroll(e.BM.Current.Buf)

	r.abuf.Append([]byte("\x1b[?25l")) // hide cursor
	r.abuf.Append([]byte("\x1b[2J"))   // clear entire screen
	r.abuf.Append([]byte("\x1b[H"))    // cursor to home

	r.drawBuffer(e.BM.Current.Buf, gutterMaxSize)
	r.drawStatusBar(e)
	r.drawCommandArea(e)

	switch e.SM.Current() {
	case string(state.Command):
		r.drawCursor(r.screenrows+2, e.CW.X()+3) // +1 to keep cursor in right spot, +2 to include the "> " prompt
	default:
		r.drawCursorOnBuffer(e.BM.Current.Buf, gutterMaxSize+1)
	}

	r.abuf.Append([]byte("\x1b[?25h")) // show the cursor
}

func (r *TextRenderer) Render(e *editor.Editor) {
	r.construct(e)
	r.abuf.Flush()
}
