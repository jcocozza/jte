package buffer

// in memory representation of the file
type Buffer interface {
	// load the file into the buffer
	Load() error
	Close() error
	// insert a character at the current cursor position
	InsertChar(c byte)
	// delete a character at the cursor
	DeleteChar()
	// insert a new line by creating a new line and moving the cursor down
	InsertNewLine()

	/* cursor stuff */

	Up() // move cursor up
	Down() // move cursor down
	Left() // move cursor left
	Right() // move cursor right
	StartLine() // move to start of line cursor is currently on
	EndLine() // move to end of line cursor is currently on
	X() int // x position of cursor
	Y() int // y position of cursor
	GoTo(x int, y int)

	// number of rows read into buffer
	NumRows() int
	// total number of rows in the file
	TotalRows() int
	Row(num int) []byte
	// if the buffer has been edited or not
	Dirty() bool
	// return the name of the buffer
	Name() string
	// convert the buffer to a single list of bytes (for saving)
	CombineRows() []byte
	// tell the buffer it has been written and dirty can be set to false
	Saved()
}

// represents a row of a file
type bufrow []byte

func (b *bufrow) InsertChar(at int, c byte) {
	if at < 0 || at > len(*b) {
		at = len(*b)
	}
	//newChars := make([]byte, len(*b)+1)
	//copy(newChars[:at], (*b)[:at])
	//newChars[at] = c
	//copy(newChars[at+1:], (*b)[at:])
	//*b = newChars
	*b = append((*b)[:at], append([]byte{c}, (*b)[at:]...)...)
}

func (b *bufrow) DelChar(at int) {
	if at < 0 || at >= len(*b) {
		return
	}
	newChars := make([]byte, len(*b)-1)
	copy(newChars[:at], (*b)[:at])
	copy(newChars[at:], (*b)[at+1:])
	*b = newChars
}

func (b *bufrow) append(bytes []byte) {
	*b = append(*b, bytes...)
}

func (b *bufrow) Trim(to int) {
	*b = (*b)[:to]
}
