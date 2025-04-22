package buffer

type FileType string

// cursor location in the buffer
//
//	X - row
//	Y - column
//
// up, down, left and right and unsafe cursor movements
// calling them directly can lead to unexpected behavior
// the cursor should only be moved at the request of its buffer,
// which knows more about what is going on and can properly limit movement
type Cursor struct {
	X int
	Y int
}

func (c *Cursor) up() {
	c.Y--
}

func (c *Cursor) down() {
	c.Y++
}

func (c *Cursor) left() {
	c.X--
}

func (c *Cursor) right() {
	c.X++
}

var SampleRows = []BufRow{
	[]byte("asdfasdfasdf"),
	[]byte("asdfasdfasdf"),
	[]byte("asdfasdfasdf"),
}

var EmptyRows = []BufRow{}

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
	FileType FileType

	// file stuff
	FileName string
	FilePath string
}

func NewBuffer(name string, readOnly bool, rows []BufRow) *Buffer {
	return &Buffer{
		Name:   name,
		Rows:   rows,
		ReadOnly: readOnly,
		cursor: &Cursor{},
	}
}

func (b *Buffer) X() int {
	return b.cursor.X
}
func (b *Buffer) Y() int {
	return b.cursor.Y
}

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
		b.cursor.up()
		b.adjustCursor()
	}
}
func (b *Buffer) Down() {
	if b.cursor.Y < len(b.Rows)-1 {
		b.cursor.down()
		b.adjustCursor()
	}
}
func (b *Buffer) Left() {
	if b.cursor.X > 0 {
		b.cursor.left()
	}
}
func (b *Buffer) Right() {
	if b.cursor.Y < len(b.Rows) && b.cursor.X < len(b.Rows[b.cursor.Y]) - 1 {
		b.cursor.right()
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
	if b.ReadOnly {return}
	if b.cursor.Y == len(b.Rows) {
		b.appendRow([]byte{})
	}
	b.Rows[b.cursor.Y].InsertChar(b.cursor.X, c)
	b.cursor.X++
	b.Modified = true
}

func (b *Buffer) DeleteChar() {
	if b.ReadOnly {return}
	if b.cursor.Y == len(b.Rows) {
		return
	}
	if b.cursor.X == 0 && b.cursor.Y == 0 {
		return
	}
	if b.cursor.X > 0 {
		b.Rows[b.cursor.Y].DelChar(b.cursor.X - 1)
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
	if b.ReadOnly {return}
	if len(b.Rows) == 0 {return}
	if len(b.Rows) == 1 {
		b.Rows[0] = BufRow(" ")
	}
	b.deleteRow(b.cursor.Y)
	b.cursor.X = 0
	b.Modified = true
}

// this the expected behavior when you press <enter>
func (b *Buffer) InsertNewLine() {
	if b.ReadOnly {return}
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
	if b.ReadOnly {return}
	if b.cursor.X == 0 {
		b.insertRow(b.cursor.Y+1, []byte(" "))
	} else {
		//b.insertRow(b.cursor.Y+1, b.Rows[b.cursor.Y][b.cursor.X:])
		b.insertRow(b.cursor.Y+1, []byte(" "))
		b.Rows[b.cursor.Y].Trim(b.cursor.X+1)
	}
	b.cursor.Y++
	b.cursor.X = 0
	b.Modified = true
}

// similar to pressing "O" in vim normal mode
func (b *Buffer) InsertNewLineAbove() {
	if b.ReadOnly {return}
	if b.cursor.X == 0 {
		b.insertRow(b.cursor.Y, []byte(" "))
	} else {
		b.insertRow(b.cursor.Y, []byte(" "))
		//b.Rows[b.cursor.Y].Trim(b.cursor.X)
	}
	b.cursor.X = 0
	b.Modified = true
}
