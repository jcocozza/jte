package buffer

import (
	"log/slog"
	"strings"

	"github.com/jcocozza/jte/internal/fileutil"
)

// an in memory representation of a file
type Buffer struct {
	logger *slog.Logger
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

	// state stuff
	Modified bool
	ReadOnly bool

	// file stuff
	FilePath string
	FileType fileutil.FileType

	CT *ChangeTracker
}

func NewBuffer(name string, filePath string, readOnly bool, rows []BufRow, l *slog.Logger) *Buffer {
	return &Buffer{
		logger: l.WithGroup("buffer"),
		Name:     name,
		FilePath: filePath,
		Rows:     rows,
		ReadOnly: readOnly,
		cursor:   &Cursor{},
		CT:       NewChangeTracker(l),
		//gutter:   &Gutter{},
	}
}

func NewEmptyBuffer(l *slog.Logger) *Buffer {
	return &Buffer{
		logger: l.WithGroup("buffer"),
		Name:     "No Name",
		ReadOnly: true,
		cursor:   &Cursor{},
		Rows:     make([]BufRow, 1),
		CT:       NewChangeTracker(l),
	}
}

func NewBufferFromString(name string, content string, l *slog.Logger) *Buffer {
	lines := strings.Split(content, "\n")

	var runes []BufRow
	for _, ln := range lines {
		if ln == "" {continue}
		l.Debug("adding line", slog.String("line", ln))
		runes = append(runes, []rune(ln))
	}

	return &Buffer{
		logger: l.WithGroup("buffer"),
		Name:     name,
		ReadOnly: true,
		cursor:   &Cursor{},
		Rows:     runes,
		CT:       NewChangeTracker(l),
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
