package buffer

import (
	"fmt"
	//"github.com/jcocozza/jte/internal/dev"
)

// cursor location in the buffer
//
//	X - column
//	Y - row
//
// up, down, left and right and unsafe cursor movements
// calling them directly can lead to unexpected behavior
// the cursor should only be moved at the request of its buffer,
// which knows more about what is going on and can properly limit movement
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

func (b *BufRow) DeleteChar(at int) {
	if at < 0 || at >= len(*b) {
		return
	}
	newChars := make([]rune, len(*b)-1)
	copy(newChars[:at], (*b)[:at])
	copy(newChars[at:], (*b)[at+1:])
	*b = newChars
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

//func (b *BufRow) append(bytes []byte) {
//	*b = append(*b, bytes...)
//}

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
}

func (b *Buffer) validCursor(cur Cursor) error {
	if cur.Y < 0 || cur.Y >= len(b.Rows) {
		return fmt.Errorf("invalid Y cursor value: %d", cur.Y)
	}
	if cur.X < 0 || cur.X > len(b.Rows[cur.Y]) {
		return fmt.Errorf("invalid X cursor value: %d", cur.X)
	}
	return nil
}

func (b *Buffer) insertRow(at int, row []rune) {
	if at < 0 || at > len(b.Rows) {
		return
	}
	b.Rows = append(b.Rows[:at], append([]BufRow{row}, b.Rows[at:]...)...)
}

// insert at a specified cursor spot
func (b *Buffer) InsertAt(at Cursor, content [][]rune) error {
	if err := b.validCursor(at); err != nil {
		//dev.Assert(err)
		return err
	}
	if len(content) >= 1 {
		err := b.Rows[at.Y].Insert(at.X, content[0])
		if err != nil {
			return err
		}
	}
	for j := 1; j < len(content); j++ {
		b.insertRow(at.Y+j, content[j])
	}
	return nil
}

// insert at the internal cursor
func (b *Buffer) Insert(content [][]rune) error {
	return b.InsertAt(*b.cursor, content)
}

// delete at a specified cursor
//
// expects that start.Y <= end.Y
// if start.Y == end.Y, then start.X < end.Y
//
// return the deleted content, empty content will be an empty list, NOT nil
func (b *Buffer) DeleteAt(start Cursor, end Cursor) ([][]rune, error) {
	if err := b.validCursor(start); err != nil {
		//dev.Assert(err)
		return nil, err
	}
	if err := b.validCursor(end); err != nil {
		//dev.Assert(err)
		return nil, err
	}
	if start.Y > end.Y {
		return nil, fmt.Errorf("invalid start/end cursors: start: %v, end: %v", start, end)
	}
	// the easy case, same line
	if start.Y == end.Y {
		content, err := b.Rows[start.Y].DeleteRange(start.X, end.X)
		if err != nil {
			return nil, err
		}
		return [][]rune{content}, nil
	}

	allDeleted := make([][]rune, end.Y-start.Y)
	// harder case, more then one line
	deletedHead, err := b.Rows[start.Y].DeleteRange(start.X, len(b.Rows[start.Y]))
	if err != nil {
		return nil, fmt.Errorf("unable to delete head from range: %w", err)
	}
	allDeleted[0] = deletedHead
	deletedTail, err := b.Rows[end.Y].DeleteRange(0, end.X+1)
	if err != nil {
		return nil, fmt.Errorf("unable to delete tail from range: %w", err)
	}
	allDeleted[len(allDeleted)-1] = deletedTail

	// splice together what is left of the start.Y and end.Y rows
	b.Rows[start.Y] = append(b.Rows[start.Y], b.Rows[end.Y]...)
	for i := start.Y + 1; i < end.Y; i++ {
		allDeleted[i-start.Y] = append([]rune(nil), b.Rows[i]...)
	}
	// actually delete everything inbetween
	b.Rows = append(b.Rows[:start.Y+1], b.Rows[end.Y+1:]...)
	return allDeleted, nil
}
