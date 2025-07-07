package buffer

// location in the buffer
//
//	X - column
//	Y - row
type Location struct {
	X int
	Y int
}

// cursor location in the buffer
type Cursor struct {
	Location
	// inclusive selection
	Selected [2]Location
}
