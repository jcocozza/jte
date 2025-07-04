package renderer

import (
	"bytes"
	"fmt"
	"log/slog"

	"github.com/jcocozza/jte/internal/buffer"
	"github.com/jcocozza/jte/internal/gutter"
)

const TAB_STOP = 8

type PaneStatusData struct {
	Active bool
	Mode string
}

type PaneRenderer interface {
	Render(rows int, cols int, psd PaneStatusData, g *gutter.Gutter, buf *buffer.Buffer) [][]byte
}

type TextPaneRenderer struct {
	rowoffset int
	coloffset int

	logger *slog.Logger
}

func NewTextPaneRenderer(l *slog.Logger) *TextPaneRenderer {
	return &TextPaneRenderer{
		logger: l.WithGroup("pane-renderer"),
	}
}

func (r *TextPaneRenderer) scroll(panerows int, panecols int, buf *buffer.Buffer) {
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

// hurestic
func runeWidth(r rune) int {
	if r < 128 {
		return 1
	}
	return 2
}

func (r *TextPaneRenderer) renderRow(row buffer.BufRow) []byte {
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

func (r *TextPaneRenderer) renderStatus(cols int, psd PaneStatusData, buf *buffer.Buffer) []byte {
	var displayModified string = ""
	if buf.Modified {
		displayModified = "(Δ)"
	}
	var displayRowNum int = 0
	totalRows := len(buf.Rows)
	currRow := buf.Y()
	if totalRows != 0 {
		displayRowNum = totalRows - 1 // -1 because i want a 0 indexed system
	}
	status := fmt.Sprintf("(%v) ln:%d/%d - %s %s", psd.Active, currRow, displayRowNum, displayModified, buf.Name)
	spacer := bytes.Repeat([]byte(" "), cols-len(status)-len(psd.Mode))
	statusBuf := append([]byte(psd.Mode), append(spacer, []byte(status)...)...)
	return statusBuf
}

func (r *TextPaneRenderer) Render(rows int, cols int, psd PaneStatusData, g *gutter.Gutter, buf *buffer.Buffer) [][]byte {
	r.scroll(rows, cols, buf)
	r.logger.Debug("rendering buffer", slog.String("name", buf.Name))
	paneBuf := make([][]byte, rows)
	for i := 0; i < rows-1; i++ {
		bufrownum := i + r.rowoffset
		if bufrownum >= len(buf.Rows) {
			paneBuf[bufrownum] = []byte("~")
			continue
		}
		paneBuf[i] = r.renderRow(buf.Rows[bufrownum])
	}
	// render status
	paneBuf[rows-1] = r.renderStatus(cols, psd, buf)
	return paneBuf
}
