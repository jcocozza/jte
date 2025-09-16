package text

import (
	"bytes"
	"log/slog"
	"strconv"

	"github.com/jcocozza/jte/internal/buffer"
)

const TAB_STOP = 8

// hurestic
func runeWidth(r rune) int {
	if r < 128 {
		return 1
	}
	return 2
}

type TextBufferRenderer struct {
	buf         *buffer.Buffer
	rowoffset   int
	coloffset   int
	gutterShift int

	logger *slog.Logger
}

func NewTextPaneRenderer(l *slog.Logger, buf *buffer.Buffer) *TextBufferRenderer {
	return &TextBufferRenderer{
		logger: l.WithGroup("pane-renderer"),
		buf:    buf,
	}
}

const scrollMargin = 10

// i hate everything about this
// it is very cursed
//
// i'm pretty sure it works though...don't mess with it
func (r *TextBufferRenderer) scroll(panerows int, panecols int, x int, y int) {
	panerows = panerows - 1 // status bar row

	// Vertical scroll
	if y < r.rowoffset+scrollMargin {
		r.rowoffset = y - scrollMargin
		if r.rowoffset < 0 {
			r.rowoffset = 0
		}
	} else if y >= r.rowoffset+panerows-scrollMargin {
		r.rowoffset = y - (panerows - scrollMargin) + 1
	}

	// Horizontal scroll
	if x < r.coloffset+scrollMargin {
		r.coloffset = x - scrollMargin
		if r.coloffset < 0 {
			r.coloffset = 0
		}
	} else if x >= r.coloffset+panecols-scrollMargin {
		r.coloffset = x - (panecols - scrollMargin) + 1
	}

	r.logger.Debug("scroll",
		slog.Int("x", x),
		slog.Int("y", y),
		slog.Int("rowoffset", r.rowoffset),
		slog.Int("coloffset", r.coloffset),
		slog.Int("panerows", panerows),
		slog.Int("panecols", panecols),
	)
}

func renderRow(row buffer.BufRow) []byte {
	var expanded []byte

	col := 0
	for _, b := range row {
		if b == '\t' {
			spaces := TAB_STOP - (col % TAB_STOP)
			expanded = append(expanded, bytes.Repeat([]byte(" "), spaces)...)
			col += spaces
		} else {
			expanded = append(expanded, []byte(string(b))...)
			col += runeWidth(b)
		}
	}
	return expanded
}

// +2 for a space on either side
func maxGutterWidth(numRows int) int {
	return len(strconv.Itoa(numRows)) + 2
}

func renderGutter(num int, maxWidth int) []byte {
	b := []byte(strconv.Itoa(num))
	repeat := max(1, maxWidth-len(b)-1)
	gutter := append(bytes.Repeat([]byte(" "), repeat), b...)
	gutter = append(gutter, []byte(" ")...)
	return gutter
}

func (r *TextBufferRenderer) render(rows int, cols int) [][]byte {
	if rows == 0 {
		rows = 1
	}
	r.scroll(rows, cols, r.buf.X(), r.buf.Y())
	paneBuf := make([][]byte, rows)
	for i := 0; i < rows-1; i++ {
		bufrownum := i + r.rowoffset
		if bufrownum >= len(r.buf.Rows) {
			paneBuf[i] = []byte("~")
			continue
		}
		maxGutterWidth := maxGutterWidth(len(r.buf.Rows))
		r.gutterShift = maxGutterWidth

		// ensure that we only render _at most_ the the number of rows
		endRow := min(cols-maxGutterWidth, len(r.buf.Rows[bufrownum]))
		if bufrownum == r.buf.Y() {
			paneBuf[i] = append(renderGutter(r.buf.Y(), maxGutterWidth), renderRow(r.buf.Rows[bufrownum][r.coloffset:endRow])...)
		} else {
			relNum := i + r.rowoffset - r.buf.Y()
			if relNum < 0 {
				relNum = relNum * -1
			}
			paneBuf[i] = append(renderGutter(relNum, maxGutterWidth), renderRow(r.buf.Rows[bufrownum][r.coloffset:endRow])...)
		}
	}
	return paneBuf
}
