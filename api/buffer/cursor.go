package buffer

// represents the cursor in the buffer
//
// X and Y cannot be less then 0
type Cursor struct {
	X int
	Y int
}

func (c *Cursor) reset() {
	c.X = 0
	c.Y = 0
}
