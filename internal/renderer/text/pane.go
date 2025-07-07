package text

import (
	"bytes"
	"log/slog"

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
	rowoffset int
	coloffset int

	logger *slog.Logger
}

func NewTextPaneRenderer(l *slog.Logger) *TextBufferRenderer {
	return &TextBufferRenderer{
		logger: l.WithGroup("pane-renderer"),
	}
}

func (r *TextBufferRenderer) scroll(panerows int, panecols int, buf *buffer.Buffer) {
	if buf.Y() < r.rowoffset {
		r.rowoffset = buf.Y()
	}
	if buf.Y() >= r.rowoffset+panerows {
		r.rowoffset = buf.Y() - panerows + 1
	}
	if buf.X() < r.coloffset {
		r.coloffset = buf.X()
	}
	if buf.X() >= r.coloffset+panecols {
		r.coloffset = buf.X() - panecols + 1
	}
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

func (r *TextBufferRenderer) render(rows int, cols int, buf *buffer.Buffer) [][]byte {
	paneBuf := make([][]byte, rows)
	for i := 0; i < rows-1; i++ {
		bufrownum := i + r.rowoffset
		if bufrownum >= len(buf.Rows) {
			paneBuf[bufrownum] = []byte("~")
			continue
		}
		paneBuf[i] = renderRow(buf.Rows[bufrownum])
	}
	// render status
	// paneBuf[rows-1] = renderStatus(cols, psd, buf)
	return paneBuf
}
