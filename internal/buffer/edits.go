package buffer

import "fmt"

// edits to the buffer

func (b *Buffer) RunningEvent() bool {
	if b.em.current == nil {
		return false
	}
	return !b.em.current.complete
}

func (b *Buffer) AcceptChange(c Change) error {
	err := c.Apply(b)
	if err != nil {
		return err
	}
	b.em.AddChange(c)
	return nil
}

func (b *Buffer) StartAndAcceptChange(c Change, etype EventType) error {
	err := c.Apply(b)
	if err != nil {
		return err
	}
	if b.em.current == nil {
		err := b.em.StartEvent(etype)
		if err != nil {
			return err
		}
	}
	b.em.AddChange(c)
	return nil
}

func (b *Buffer) Commit() {
	b.em.Commit()
}

func (b *Buffer) StartEvent(etype EventType) {
	b.em.StartEvent(etype)
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

func (b *Buffer) deleteRow(at int) ([]rune, error) {
	if at < 0 || at >= len(b.Rows) {
		return  nil, fmt.Errorf("cannot delete row at %d", at)
	}
	content := b.Rows[at]
	b.Rows = append(b.Rows[:at], b.Rows[at+1:]...)
	return content, nil
}

// insert at a specified cursor spot
func (b *Buffer) insertAt(at Cursor, content [][]rune) error {
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
func (b *Buffer) insert(content [][]rune) error {
	return b.insertAt(*b.cursor, content)
}

// delete at a specified cursor
//
// expects that start.Y <= end.Y
// if start.Y == end.Y, then start.X < end.Y
//
// return the deleted content, empty content will be an empty list, NOT nil
func (b *Buffer) deleteAt(start Cursor, end Cursor) ([][]rune, error) {
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

// wrapper around DeleteAt for the current cursor position
func (b *Buffer) delete() ([][]rune, error) {
	return b.deleteAt(*b.cursor, *b.cursor)
}
