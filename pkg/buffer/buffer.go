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
	if b.cursor.Y < len(b.Rows) && b.cursor.X < len(b.Rows[b.cursor.Y]) {
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
