package actions

type Action int

const (
	None Action = iota
	Exit
	Repeat
	CursorUp
	CursorDown
	CursorRight
	CursorLeft
)
