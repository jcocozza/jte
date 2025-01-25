package buffer

import (
	"bufio"
	"bytes"
	"fmt"

	"github.com/jcocozza/jte/api/fileutil"
)

type bufrow []byte

// in memory representation of the file
type Buffer struct {
	Rows []*bufrow
	C    *Cursor

	Name string
	Dirty bool
}

func NewEmptyBuffer() *Buffer {
	return &Buffer{
		Rows: []*bufrow{},
		C: &Cursor{},
	}
}

func (b *Buffer) IsEmpty() bool {
	return len(b.Rows) == 0
}

func (b *Buffer) insertRow(at int, row []byte) {
	if at < 0 || at > len(b.Rows) {
		return
	}
	newRows := make([]*bufrow, len(b.Rows)+1)
	copy(newRows[:at], b.Rows[:at])
	newBufRow := bufrow(row)
	newRows[at] = &newBufRow
	copy(newRows[at+1:], b.Rows[at:])
	b.Rows = newRows
}

func (b *Buffer) LoadFromBytes(rows [][]byte) {
	for i, row := range rows {
		b.insertRow(i, row)
	}
}

// read the file into the buffer
func (b *Buffer) Load(filename string) error {
	file, err := fileutil.OpenOrCreateFile(filename)
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(file)
	numScans := 0
	for scanner.Scan() {
		line := scanner.Bytes()
		line = bytes.TrimRight(line, "\r\n")
		b.insertRow(len(b.Rows), line)
		numScans += 1
	}
	if numScans == 0 {
		b.insertRow(0, []byte{})
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}
	b.Name = filename
	return nil
}

// when moving up or down and at the end of a line, we want to snap to end of next line if that line is shorter
func (b *Buffer) adjustCursor() {
	if b.C.Y >= len(b.Rows) {
		return
	}
	newRowLen := len(*b.Rows[b.C.Y])
	if b.C.X > newRowLen {
		b.C.X = newRowLen
	}
}
func (b *Buffer) Up() {
	if b.C.Y > 0 {
		b.C.up()
		b.adjustCursor()
	}
}
func (b *Buffer) Down() {
	if b.C.Y < len(b.Rows) -1 {
		b.C.down()
		b.adjustCursor()
	}
}
func (b *Buffer) Left() {
	if b.C.X > 0 {
		b.C.left()
	}
}
func (b *Buffer) Right() {
	if b.C.Y < len(b.Rows) && b.C.X < len(*b.Rows[b.C.Y]) {
		b.C.right()
	}
}
