package buffer

// a change is something applied to a buffer
type Change interface {
	Apply(buf *Buffer) error
}

type InsertAt struct {
	cur      Cursor
	contents [][]rune
}

func (i InsertAt) Apply(buf *Buffer) error {
	return buf.insertAt(i.cur, i.contents)
}

// insert at the buffer's internal cursor
type Insert struct {
	Contents [][]rune
}

func (i Insert) Apply(buf *Buffer) error {
	return buf.insert(i.Contents)
}

// insert new line at buffer's interal cursor
type EnterNewLine struct {}
func (i EnterNewLine) Apply(buf *Buffer) error {
	buf.insertRow([]rune{})
	return nil
}

type DeleteAt struct {
	StartCur, EndCur Cursor
	Contents         [][]rune
}

func (d DeleteAt) Apply(buf *Buffer) error {
	contents, err := buf.deleteAt(d.StartCur, d.EndCur)
	if err != nil {
		return err
	}
	d.Contents = contents
	return nil
}

// delete at buffer's interal cursor
type Delete struct {
	contents [][]rune
}

func (d Delete) Apply(buf *Buffer) error {
	contents, err := buf.delete()
	if err != nil {
		return err
	}
	d.contents = contents
	return nil
}

// backspace at buffer's internal cursor
type Backspace struct {
	contents [][]rune
}

func (b Backspace) Apply(buf *Buffer) error {
	content, err := buf.backspace()
	if err != nil {
		return err
	}
	b.contents = content
	return nil
}


// delete line at the cursor
type DeleteLine struct{ contents []rune }

func (d DeleteLine) Apply(buf *Buffer) error {
	content, err := buf.deleteRow(buf.cursor.Y)
	if err != nil {
		return err
	}
	d.contents = content
	return nil
}
