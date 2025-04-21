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

	InsertChar

	Mode_Normal
	Mode_Insert
)

var ActionNames = [...]string{
	None: "none",
	Exit: "exit",
	Repeat: "repeat",
	CursorUp: "cursor up",
	CursorDown: "cursor down",
	CursorLeft: "cursor left",
	CursorRight: "cursor right",

	InsertChar: "insert char",

	Mode_Normal: "mode normal",
	Mode_Insert: "mode insert",
}
