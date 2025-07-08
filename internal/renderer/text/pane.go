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
	buf *buffer.Buffer
	rowoffset int
	coloffset int
	gutterShift int

	logger *slog.Logger
}

func NewTextPaneRenderer(l *slog.Logger, buf *buffer.Buffer) *TextBufferRenderer {
	return &TextBufferRenderer{
		logger: l.WithGroup("pane-renderer"),
		buf: buf,
	}
}

func (r *TextBufferRenderer) scroll(panerows int, panecols int, x int, y int) {
	if y < r.rowoffset {
		r.rowoffset = y
	}
	panerows = panerows - 1 // leave room for the status bar in each pane
	if y >= r.rowoffset+panerows {
		r.rowoffset = y - panerows + 1
	}
	if x < r.coloffset {
		r.coloffset = x
	}
	if x >= r.coloffset+panecols {
		r.coloffset = x - panecols + 1
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
	repeat := maxWidth - len(b) - 1
	if repeat < 1 {
		repeat = 1
	}
	gutter := append(bytes.Repeat([]byte(" "), repeat), b...)
	gutter = append(gutter, []byte(" ")...)
	return gutter
}

func (r *TextBufferRenderer) render(rows int, cols int) [][]byte {
	r.scroll(rows, cols, r.buf.X(), r.buf.Y())
	paneBuf := make([][]byte, rows)
	for i := 0; i < rows-1; i++ {
		bufrownum := i + r.rowoffset
		if bufrownum >= len(r.buf.Rows) {
			paneBuf[bufrownum] = []byte("~")
			continue
		}
		//r.logger.Debug("rendering row", slog.Int("bufrownum", bufrownum), slog.Int("coloffset", r.coloffset))

		maxGutterWidth := maxGutterWidth(len(r.buf.Rows))
		r.gutterShift = maxGutterWidth
		if bufrownum == r.buf.Y() {
			paneBuf[i] = append(renderGutter(r.buf.Y(), maxGutterWidth), renderRow(r.buf.Rows[bufrownum][r.coloffset:])...)
		} else {
			relNum := i + r.rowoffset - r.buf.Y()
			if relNum < 0 {
				relNum = relNum * -1
			}
			paneBuf[i] = append(renderGutter(relNum, maxGutterWidth), renderRow(r.buf.Rows[bufrownum][r.coloffset:])...)
		}
	}
	return paneBuf
}
