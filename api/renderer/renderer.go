package renderer

import (
	"bytes"
	"fmt"
	"log/slog"
	"os"

	"github.com/jcocozza/jte/api/buffer"
	"github.com/jcocozza/jte/api/command"
	"github.com/jcocozza/jte/term"
)

const TAB_STOP = 8

type Renderer interface {
	Init(l *slog.Logger) error
	Render(buf buffer.Buffer, s StatusInfo, cw *command.CommandWindow)
	Cleanup()
	Exit(msg string)
	ExitErr(err error)
	//SetMsg(s StatusInfo, buf buffer.Buffer, msg messages.Message)
}

type TextRenderer struct {
	screenrows int
	initscreenrows int
	screencols int

	rowoffset int
	coloffset int

	rw *term.RawMode

	abuf abuf

	//currMsg messages.Message

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
	r.screenrows = rows - 1 // leave room for status bar
	r.initscreenrows = rows - 1 // leave room for status bar
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

func (r *TextRenderer) drawCursor(buf buffer.Buffer) {
	y := (buf.Y() - r.rowoffset) + 1
	x := (buf.X() - r.coloffset) + 1
	s := fmt.Sprintf("\x1b[%d;%dH", y, x)
	r.abuf.Append([]byte(s))
}

//func (r *TextRenderer) SetMsg(s StatusInfo, buf buffer.Buffer, msg messages.Message, cw command.CommandWindow) {
//	r.currMsg = msg
//	r.Render(buf, s, cw)
//}

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

type StatusInfo struct {
	Name      string
	Dirty     bool
	CurrRow   int
	TotalRows int
	Mode      string
}

func (s *StatusInfo) String() string {
	var displayName string = s.Name
	if displayName == "" {
		displayName = "[No Name]"
	}
	var displayDirty string = ""
	if s.Dirty {
		displayDirty = "(modified)"
	}
	var displayRowNum int = 0
	if s.TotalRows != 0 {
		displayRowNum = s.TotalRows - 1 // -1 because i want a 0 indexed system
	}
	return fmt.Sprintf("ln:%d/%d - %s %s", s.CurrRow, displayRowNum, displayDirty, displayName)
}

func (r *TextRenderer) drawStatusBar(s StatusInfo, cw *command.CommandWindow) {
	// status bar
	r.abuf.Append([]byte("\x1b[7m"))
	mode := s.Mode
	status := s.String()

	r.abuf.Append([]byte(mode))
	r.abuf.Append(bytes.Repeat([]byte(" "), r.screencols-len(status)-len(mode)))
	r.abuf.Append([]byte(status))
	r.abuf.Append([]byte("\x1b[m"))
	r.abuf.Append([]byte("\r\n"))

	r.abuf.Append([]byte("\x1b[K"))
	o := cw.Output()
	if len(o) > 0 {
		for i, row := range o {
			r.abuf.Append(row)
			if i < len(o) {
				r.abuf.Append([]byte("\r\n"))
			}
		}
	}
	r.abuf.Append(cw.Prompt())
	// messages
	//var displayMsg string = r.currMsg.Text
	//if len(r.currMsg.Text) > r.screencols {
	//	displayMsg = r.currMsg.Text[0:r.screencols]
	//}
	//if r.currMsg.NonEmpty() && !r.currMsg.Expired() {
	//	r.abuf.Append([]byte(displayMsg))
	//}
}

func (r *TextRenderer) renderRow(row []byte) {
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
	r.abuf.Append(expanded)
}

func (r *TextRenderer) drawBuffer(buf buffer.Buffer) {
	for i := 0; i < r.screenrows; i++ {
		filerow := i + r.rowoffset
		if filerow >= buf.NumRows() {
			r.abuf.Append([]byte("~"))
		} else {
			//r.logger.Debug("rendering", slog.String("row", string(buf.Row(filerow))), slog.Int("num", filerow))
			//r.abuf.Append(buf.Row(filerow))
			r.renderRow(buf.Row(filerow))
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

func (r *TextRenderer) Render(buf buffer.Buffer, s StatusInfo, cw *command.CommandWindow) {
	r.screenrows = r.initscreenrows - cw.NumRows()
	r.scroll(buf)
	r.abuf.Append([]byte("\x1b[?25l")) // hide cursor
	r.abuf.Append([]byte("\x1b[2J")) // clear entire screen
	r.abuf.Append([]byte("\x1b[H")) // cursor to home

	if buf.NumRows() == 0 {
		r.drawWelcome()
	} else {
		r.drawBuffer(buf)
	}
	r.drawStatusBar(s, cw)

	r.drawCursor(buf)
	r.abuf.Append([]byte("\x1b[?25h"))
	r.abuf.Flush()
}
