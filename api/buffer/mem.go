package buffer

import (
	"bufio"
	"bytes"
	"fmt"
	"log/slog"

	"github.com/jcocozza/jte/api/fileutil"
)

type MemBuffer struct {
	Rows []*bufrow
	C    *Cursor

	name  string
	dirty bool

	logger *slog.Logger
}

func NewEmptyBuffer(name string, l *slog.Logger) *MemBuffer {
	return &MemBuffer{
		Rows:   []*bufrow{},
		C:      &Cursor{},
		name:   name,
		logger: l,
	}
}

// appending is faster then inserting
func (b *MemBuffer) appendRow(row []byte) {
	b.logger.Debug("append row", slog.String("row", string(row)), slog.Int("row", len(b.Rows)))
	newBufRow := make(bufrow, len(row))
	copy(newBufRow, row) // Ensure a copy of the row to avoid unintended aliasing
	b.Rows = append(b.Rows, &newBufRow)
}

func (b *MemBuffer) insertRow(at int, row []byte) {
	if at < 0 || at > len(b.Rows) {
		return
	}
	b.logger.Debug("insert row", slog.String("row", string(row)), slog.Int("at", at))
	newRows := make([]*bufrow, len(b.Rows)+1)
	copy(newRows[:at], b.Rows[:at])
	newBufRow := make(bufrow, len(row))
	copy(newBufRow, row) // Ensure a copy of the row to avoid unintended aliasing
	newRows[at] = &newBufRow
	copy(newRows[at+1:], b.Rows[at:])
	b.Rows = newRows
}

func (b *MemBuffer) LoadFromBytes(rows [][]byte) {
	for i, row := range rows {
		b.insertRow(i, row)
	}
}

// read the file into the buffer
func (b *MemBuffer) Load() error {
	file, err := fileutil.OpenOrCreateFile(b.name)
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(file)
	numScans := 0
	for scanner.Scan() {
		line := scanner.Bytes()
		line = bytes.TrimRight(line, "\r\n")
		b.appendRow(line)
		numScans += 1
	}
	if numScans == 0 {
		b.insertRow(0, []byte{})
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}
	//b.Name = filename
	return nil
}

// when moving up or down and at the end of a line, we want to snap to end of next line if that line is shorter
func (b *MemBuffer) adjustCursor() {
	if b.C.Y >= len(b.Rows) {
		return
	}
	newRowLen := len(*b.Rows[b.C.Y])
	if b.C.X > newRowLen {
		b.C.X = newRowLen
	}
}
func (b *MemBuffer) Up() {
	if b.C.Y > 0 {
		b.C.up()
		b.adjustCursor()
	}
}
func (b *MemBuffer) Down() {
	if b.C.Y < len(b.Rows)-1 {
		b.C.down()
		b.adjustCursor()
	}
}
func (b *MemBuffer) Left() {
	if b.C.X > 0 {
		b.C.left()
	}
}
func (b *MemBuffer) Right() {
	if b.C.Y < len(b.Rows) && b.C.X < len(*b.Rows[b.C.Y]) {
		b.C.right()
	}
}
func (b *MemBuffer) X() int {
	return b.C.X
}
func (b *MemBuffer) Y() int {
	return b.C.Y
}
func (b *MemBuffer) NumRows() int {
	return len(b.Rows)
}
func (b *MemBuffer) TotalRows() int {
	return len(b.Rows)
}
func (b *MemBuffer) Row(num int) []byte {
	return *b.Rows[num]
}
func (b *MemBuffer) Dirty() bool {
	return b.dirty
}
func (b *MemBuffer) Name() string {
	return b.name
}
