package actions

type Action int

const (
	None Action = iota
	Exit
	Repeat

	// Navigation

	CursorUp
	CursorDown
	CursorRight
	CursorLeft
	StartLine
	EndLine

	// editing

	InsertChar
	InsertNewLine
	InsertNewLineBelow
	InsertNewLineAbove
	DeleteChar
	RemoveChar
	DeleteLine

	// modality

	Mode_Normal
	Mode_Insert
)

// this this just a convience for debugging purposes
//
// allows us to easily map the "enum" of actions to a string representation
var ActionNames = [...]string{
	None:   "none",
	Exit:   "exit",
	Repeat: "repeat",

	CursorUp:    "cursor up",
	CursorDown:  "cursor down",
	CursorLeft:  "cursor left",
	CursorRight: "cursor right",
	StartLine:   "start line",
	EndLine:     "end line",

	InsertChar:         "insert char",
	InsertNewLine:      "insert new line",
	InsertNewLineBelow: "insert new line below",
	InsertNewLineAbove: "insert new line above",
	DeleteChar:         "delete char",
	RemoveChar:         "remove char",
	DeleteLine:         "delete line",

	Mode_Normal: "mode normal",
	Mode_Insert: "mode insert",
}
