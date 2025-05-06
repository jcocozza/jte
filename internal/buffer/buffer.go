package buffer

import (
	"fmt"
	"log/slog"

	"github.com/jcocozza/jte/internal/fileutil"
	//"github.com/jcocozza/jte/internal/dev"
)

// cursor location in the buffer
//
//	X - column
//	Y - row
type Cursor struct {
	X int
	Y int
}

// Represents a single row in the buffer
type BufRow []rune

func (b *BufRow) Insert(at int, content []rune) error {
	if at < 0 || at > len(*b) {
		return fmt.Errorf("invalid row value: %d", at)
	}
	// i think this is a trick to do a little less work when just appending to the end of the row
	if at == len(*b) {
		*b = append(*b, content...)
		return nil
	}
	*b = append((*b)[:at], append(content, (*b)[at:]...)...)
	return nil
}

func (b *BufRow) DeleteChar(at int) (rune, error) {
	if at < 0 || at >= len(*b) {
		return -1, fmt.Errorf("invalid loc to delete")
	}
	char := (*b)[at]
	newChars := make([]rune, len(*b)-1)
	copy(newChars[:at], (*b)[:at])
	copy(newChars[at:], (*b)[at+1:])
	*b = newChars
	return char, nil
}

func (b *BufRow) DeleteRange(start, end int) ([]rune, error) {
	if start < 0 || start > len(*b) {
		return nil, fmt.Errorf("invalid start index: %d", start)
	}
	if end < 0 || end > len(*b) {
		return nil, fmt.Errorf("invalid end index: %d", end)
	}
	if start > end {
		return nil, fmt.Errorf("start cannot be greater than end: %d > %d", start, end)
	}
	content := append([]rune(nil), (*b)[start:end]...)
	*b = append((*b)[:start], (*b)[end:]...)
	return content, nil
}

func (b *BufRow) append(runes []rune) {
	*b = append(*b, runes...)
}

//func (b *BufRow) Trim(to int) {
//	*b = (*b)[:to]
//}

// an in memory representation of a file
type Buffer struct {
	// a unique identifier
	id int
	// purely for display purposes
	// in most cases, this will be the same as the file name
	// however, sometimes we just want a quick buffer
	// in this case, we use another name
	Name string

	// the rows in the underlying file
	Rows   []BufRow
	cursor *Cursor

	// state stuff
	Modified bool
	ReadOnly bool

	// file stuff
	FilePath string
	FileType fileutil.FileType

	// events
	em *EventManager
}

func NewBuffer(name string, filePath string, readOnly bool, rows []BufRow, l *slog.Logger) *Buffer {
	return &Buffer{
		Name:     name,
		FilePath: filePath,
		Rows:     rows,
		ReadOnly: readOnly,
		cursor:   &Cursor{},
		em:       NewEventManager(l),
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
