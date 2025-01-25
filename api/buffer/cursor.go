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

// these are bare cursor movements.
// they are NOT safe and can lead to unexpected behavior.
// the cursor should only be moved at the request of its buffer,
// which know more about what is going on and can properly limit movement

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
