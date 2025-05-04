package buffer

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

// Represents a single row in the buffer
type BufRow []byte

func (b *BufRow) InsertChar(at int, c byte) {
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

func (b *BufRow) DelChar(at int) {
	if at < 0 || at >= len(*b) {
		return
	}
	newChars := make([]byte, len(*b)-1)
	copy(newChars[:at], (*b)[:at])
	copy(newChars[at:], (*b)[at+1:])
	*b = newChars
}

func (b *BufRow) append(bytes []byte) {
	*b = append(*b, bytes...)
}

func (b *BufRow) Trim(to int) {
	*b = (*b)[:to]
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

	// file stuff
	FilePath string
}
