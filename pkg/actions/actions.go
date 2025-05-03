package actions

type Action int

// steps for adding an action:
// 1. create a new const
// 2. map to an action name (add to ActionNames)
// 3. add action to registry (pkg/editor/registry.go)
// optional steps (you will likely do one or the other):
// - use action in relevant bindings (pkg/bindings/default.go)
// - use action in the command registry (pkg/commandWindow/commands.go)

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
	Mode_Command

	// command mode stuff

	InsertCommandChar
	DeleteCommandChar
	ClearCommand
	Submit
	RunCommand

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

	Mode_Normal:  "mode normal",
	Mode_Insert:  "mode insert",
	Mode_Command: "mode command",

	InsertCommandChar: "insert command char",
	DeleteCommandChar: "delete command char",
	ClearCommand:      "clear command",
	Submit:            "submit",

	RunCommand: "run command",
}
