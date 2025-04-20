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


// represents a single row
type row []byte

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
	rows   []row
	cursor Cursor

	// state stuff
	Modified bool
	ReadOnly bool
	FileType FileType

	// file stuff
	FileName string
	FilePath string
}
