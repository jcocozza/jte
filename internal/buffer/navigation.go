package buffer

// navigation in the buffer

// when moving up or down and at the end of a line, we want to snap to end of next line if that line is shorter
func (b *Buffer) adjustCursor() {
	if b.cursor.Y >= len(b.Rows) {
		return
	}
	newRowLen := len(b.Rows[b.cursor.Y])
	if b.cursor.X > newRowLen {
		b.cursor.X = newRowLen
	}
}

func (b *Buffer) Up() {
	if b.cursor.Y > 0 {
		b.cursor.Y--
		b.adjustCursor()
	}
}
func (b *Buffer) Down() {
	if b.cursor.Y < len(b.Rows)-1 {
		b.cursor.Y++
		b.adjustCursor()
	}
}
func (b *Buffer) Left() {
	if b.cursor.X > 0 {
		b.cursor.X--
	}
}
func (b *Buffer) Right() {
	if b.cursor.Y < len(b.Rows) && b.cursor.X < len(b.Rows[b.cursor.Y])-1 {
		b.cursor.X++
	}
}
