package bindings

type ActionId int

const (
	Exit ActionId = iota

	Action_Commit

	// Modality

	Action_ModeNormal
	Action_ModeInsert
	Action_ModeCommand

	// Navigation

	Action_CursorUp
	Action_CursorDown
	Action_CursorRight
	Action_CursorLeft
	Action_StartLine
	Action_EndLine

	// editing

	Action_InsertChar
	Action_InsertNewLine
	Action_InsertNewLineBelow
	Action_InsertNewLineAbove
	Action_DeleteChar
	Action_RemoveChar
	Action_DeleteLine

	// command mode stuff

	InsertCommandChar
	DeleteCommandChar
	ClearCommand
	// Submit
	// RunCommand
)
