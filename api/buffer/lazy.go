package buffer

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
)

func countLines(file *os.File) (int, error) {
	scanner := bufio.NewScanner(file)
	num := 0
	for scanner.Scan() {
		num++
	}
	if err := scanner.Err(); err != nil {
		return 0, err
	}
	return num, nil
}

type LazyBuffer struct {
	Rows []*bufrow
	C    *Cursor

	name        string
	File        *os.File
	dirty       bool
	LoadedLines int
	BufSize     int
	totalLines  int

	logger *slog.Logger
}

func NewLazyBuffer(filename string, bufSize int, l *slog.Logger) (*LazyBuffer, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	totalLines, err := countLines(file)
	if err != nil {
		return nil, err
	}
	_, err = file.Seek(0, 0)
	if err != nil {
		return nil, fmt.Errorf("error resetting file position: %w", err)
	}
	return &LazyBuffer{
		Rows:       []*bufrow{},
		C:          &Cursor{},
		name:       filename,
		File:       file,
		BufSize:    bufSize,
		totalLines: totalLines,
		logger: l,
	}, nil
}

func (b *LazyBuffer) LoadNextChunk() error {
	if b.LoadedLines >= b.totalLines {
		return nil
	}
	startLine := b.LoadedLines
	endLine := startLine + b.BufSize
	if endLine > b.totalLines {
		endLine = b.totalLines
	}

	// Seek to the beginning of the file
	_, err := b.File.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("error seeking file: %w", err)
	}

	scanner := bufio.NewScanner(b.File)
	currentLine := 0
	for scanner.Scan() {
		// Skip lines until startLine
		if currentLine < startLine {
			currentLine++
			continue
		}
		// Process lines in the desired range
		if currentLine >= startLine && currentLine < endLine {
			b.insertRow(currentLine-startLine, scanner.Bytes())
		}
		currentLine++
		if currentLine >= endLine {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	b.LoadedLines = endLine
	return nil
}

// appending is faster then inserting
func (b *LazyBuffer) appendRow(row []byte) {
	b.logger.Debug("append row", slog.String("row", string(row)), slog.Int("row", len(b.Rows)))
	newBufRow := make(bufrow, len(row))
	copy(newBufRow, row) // Ensure a copy of the row to avoid unintended aliasing
	b.Rows = append(b.Rows, &newBufRow)
}

func (b *LazyBuffer) insertRow(at int, row []byte) {
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

func (b *LazyBuffer) adjustCursor() {
	if b.C.Y >= len(b.Rows) {
		return
	}
	newRowLen := len(*b.Rows[b.C.Y])
	if b.C.X > newRowLen {
		b.C.X = newRowLen
	}
}

// Load loads the file and starts lazy loading
func (b *LazyBuffer) Load() error {
	// Load the first chunk
	return b.LoadNextChunk()
}

// Up moves the cursor up, loading more lines if necessary
func (b *LazyBuffer) Up() {
	if b.C.Y > 0 {
		b.C.up()
		b.adjustCursor()
	}
}

// Down moves the cursor down, loading more lines if necessary
func (b *LazyBuffer) Down() {
	if b.C.Y < len(b.Rows)-1 {
		b.C.down()
		b.adjustCursor()
	} else {
		// If we're at the bottom of the buffer, load the next chunk
		if b.LoadedLines < b.totalLines {
			err := b.LoadNextChunk()
			if err != nil {
				panic(err)
			}
		}
	}
}

// Left moves the cursor left
func (b *LazyBuffer) Left() {
	if b.C.X > 0 {
		b.C.left()
	}
}

// Right moves the cursor right
func (b *LazyBuffer) Right() {
	if b.C.Y < len(b.Rows) && b.C.X < len(*b.Rows[b.C.Y]) {
		b.C.right()
	}
}

func (b *LazyBuffer) X() int {
	return b.C.X
}
func (b *LazyBuffer) Y() int {
	return b.C.Y
}
func (b *LazyBuffer) NumRows() int {
	return len(b.Rows)
}
func (b *LazyBuffer) TotalRows() int {
	return b.totalLines
}
func (b *LazyBuffer) Row(num int) []byte {
	return *b.Rows[num]
}
func (b *LazyBuffer) Dirty() bool {
	return b.dirty
}
func (b *LazyBuffer) Name() string {
	return b.name
}
