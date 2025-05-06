package gutter

type GutterRow struct {
	Num int
	Content []rune
}

type Gutter struct {
	Rows []GutterRow
}

func NewGutter() *Gutter {
	return &Gutter{
		Rows: []GutterRow{},
	}
}
