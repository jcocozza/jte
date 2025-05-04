package buffer

// these are things that can happen in the buffer

// NAVIGATION
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

// go to start of current line
func (b *Buffer) StartLine() {
	b.cursor.X = 0
}

// go to end of current line
func (b *Buffer) EndLine() {
	if b.cursor.Y < len(b.Rows) {
		b.cursor.X = len(b.Rows[b.cursor.Y])
	}
}

// EDITING

func (b *Buffer) insertRow(at int, row []byte) {
	if at < 0 || at > len(b.Rows) {
		return
	}
	b.Rows = append(b.Rows[:at], append([]BufRow{row}, b.Rows[at:]...)...)
}

func (b *Buffer) appendRow(row []byte) {
	newBufRow := make(BufRow, len(row))
	copy(newBufRow, row) // Ensure a copy of the row to avoid unintended aliasing
	b.Rows = append(b.Rows, newBufRow)
}

func (b *Buffer) deleteRow(at int) {
	if at < 0 || at >= len(b.Rows) {
		return
	}
	b.Rows = append(b.Rows[:at], b.Rows[at+1:]...)
}

func (b *Buffer) InsertChar(c byte) {
	if b.ReadOnly {
		return
	}
	if b.cursor.Y == len(b.Rows) {
		b.appendRow([]byte{})
	}
	b.Rows[b.cursor.Y].InsertChar(b.cursor.X, c)
	b.cursor.X++
	b.Modified = true
}

func (b *Buffer) DeleteChar() {
	if b.ReadOnly {
		return
	}
	if b.cursor.Y == len(b.Rows) {
		return
	}
	if b.cursor.X == 0 && b.cursor.Y == 0 {
		return
	}
	if b.cursor.X > 0 {
		b.Rows[b.cursor.Y].DeleteChar(b.cursor.X - 1)
		b.cursor.X--
	} else {
		newX := len(b.Rows[b.cursor.Y-1])
		b.Rows[b.cursor.Y-1].append(b.Rows[b.cursor.Y])
		b.deleteRow(b.cursor.Y)
		b.cursor.Y--
		b.cursor.X = newX
	}
	b.Modified = true
}

func (b *Buffer) DeleteLine() {
	if b.ReadOnly {
		return
	}
	if len(b.Rows) == 0 {
		return
	}
	if len(b.Rows) == 1 {
		b.Rows[0] = BufRow(" ")
	}
	b.deleteRow(b.cursor.Y)
	b.cursor.X = 0
	b.Modified = true
}

// this the expected behavior when you press <enter>
func (b *Buffer) InsertNewLine() {
	if b.ReadOnly {
		return
	}
	if b.cursor.X == 0 {
		b.insertRow(b.cursor.Y, []byte(" "))
	} else {
		b.insertRow(b.cursor.Y+1, b.Rows[b.cursor.Y][b.cursor.X:])
		b.Rows[b.cursor.Y].Trim(b.cursor.X)
	}
	b.cursor.Y++
	b.cursor.X = 0
	b.Modified = true
}

// similar to pressing "o" in vim normal mode
func (b *Buffer) InsertNewLineBelow() {
	if b.ReadOnly {
		return
	}
	if b.cursor.X == 0 {
		b.insertRow(b.cursor.Y+1, []byte(" "))
	} else {
		//b.insertRow(b.cursor.Y+1, b.Rows[b.cursor.Y][b.cursor.X:])
		b.insertRow(b.cursor.Y+1, []byte(" "))
		b.Rows[b.cursor.Y].Trim(b.cursor.X + 1)
	}
	b.cursor.Y++
	b.cursor.X = 0
	b.Modified = true
}

// similar to pressing "O" in vim normal mode
func (b *Buffer) InsertNewLineAbove() {
	if b.ReadOnly {
		return
	}
	if b.cursor.X == 0 {
		b.insertRow(b.cursor.Y, []byte(" "))
	} else {
		b.insertRow(b.cursor.Y, []byte(" "))
		//b.Rows[b.cursor.Y].Trim(b.cursor.X)
	}
	b.cursor.X = 0
	b.Modified = true
}
