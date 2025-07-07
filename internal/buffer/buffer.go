package buffer

import (
	"log/slog"

	"github.com/jcocozza/jte/internal/fileutil"
)

// an in memory representation of a file
type Buffer struct {
	// a unique identifier
	id int
	// purely for display purposes
	// in most cases, this will be the same as the file name
	// however, sometimes we just want a quick buffer
	// in this case, we use another name
	Name string

	// the actual data
	Rows []BufRow

	// cursors in the buffer
	// main cursor
	cursor *Cursor
	// extra cursors
	cursors []*Cursor

	gutter *Gutter

	// state stuff
	Modified bool
	ReadOnly bool

	// file stuff
	FilePath string
	FileType fileutil.FileType
}

func NewBuffer(name string, filePath string, readOnly bool, rows []BufRow, l *slog.Logger) *Buffer {
	return &Buffer{
		Name:     name,
		FilePath: filePath,
		Rows:     rows,
		ReadOnly: readOnly,
		cursor:   &Cursor{},
		gutter:   &Gutter{},
	}
}

func ReadFileIntoBuffer(path string, l *slog.Logger) (*Buffer, error) {
	content, writeable, ftype, err := fileutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	readOnly := !writeable
	bufrows := make([]BufRow, len(content))
	for i, row := range content {
		bufrows[i] = BufRow(row)
	}
	buf := NewBuffer(path, path, readOnly, bufrows, l)
	buf.FileType = ftype
	return buf, nil
}
