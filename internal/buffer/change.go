package buffer

import "fmt"

type Change interface {
	// for debugging
	String() string
	Do(b *Buffer)
	Undo(b *Buffer)
}

type Insert struct {
	Loc  Location
	Char rune
}

func (c *Insert) String() string {
	return fmt.Sprintf("insert: %s at (%d,%d)", string(c.Char), c.Loc.X, c.Loc.Y)
}

func (c *Insert) Do(b *Buffer) {
	b.Insert(c.Char, c.Loc)
	//b.cursor.Location = c.Loc
}

func (c *Insert) Undo(b *Buffer) {
	b.Delete(c.Loc)
}

type Delete struct {
	Loc  Location
	Char *rune // nil before deletion takes place
}

func (c *Delete) String() string {
	return fmt.Sprintf("delete: at (%d,%d)", c.Loc.X, c.Loc.Y)
}

func (c *Delete) Do(b *Buffer) {
	char := b.Delete(c.Loc)
	c.Char = &char
}

func (c *Delete) Undo(b *Buffer) {
	b.Insert(*c.Char, c.Loc)
}

// a set of changes that are done and undone together
type ChangeBlock struct {
	block []Change
	locked bool
}

// do all changes in the block
func (cb *ChangeBlock) Do(b *Buffer) {
	for _, c := range cb.block {
		c.Do(b)
	}
}

// undo all changes in the block
func (cb *ChangeBlock) Undo(b *Buffer) {
	// because block is created sequentially, we need to undo in reverse order of block
	for i := len(cb.block) - 1; i >= 0; i-- {
		cb.block[i].Undo(b)
		//c.Undo(b)
	}
}

// TODO: trying to add changed after lock should raise alarm bells
func (cb *ChangeBlock) Add(c Change) {
	if cb.locked { return }
	cb.block = append(cb.block, c)
}

// when locked no new changes should be added to the block
func (cb *ChangeBlock) Lock() {
	cb.locked = true
}
