package buffer

type GutterRow struct {
	Num int
	Content []rune
}

type Gutter struct {
	Rows []GutterRow
}
