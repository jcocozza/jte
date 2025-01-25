package buffer

// in memory representation of the file
type Buffer interface {
	// load the file into the buffer
	Load() error
	// move cursor up
	Up()
	// move cursor down
	Down()
	// move cursor left
	Left()
	// move cursor right
	Right()
	// X position of the cursor in the buffer
	X() int
	// Y position of the cursor in the buffer
	Y() int
	// number of rows read into buffer
	NumRows() int
	// total number of rows in the file
	TotalRows() int
	Row(num int) []byte
	Dirty() bool
	// return the name of the buffer
	Name() string
}

// represents a row of a file
type bufrow []byte

