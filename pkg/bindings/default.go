package bindings

import (
	"github.com/jcocozza/jte/pkg/actions"
	"github.com/jcocozza/jte/pkg/keyboard"
)

var Normal = &BindingNode{
	actions: []actions.Action{},
	children: map[keyboard.Key]*BindingNode{
		keyboard.CtrlC: {children: nil, actions: []actions.Action{actions.Exit}},
		'i':            {children: nil, actions: []actions.Action{actions.Mode_Insert}},

		'h': {children: nil, actions: []actions.Action{actions.CursorLeft}},
		'j': {children: nil, actions: []actions.Action{actions.CursorDown}},
		'k': {children: nil, actions: []actions.Action{actions.CursorUp}},
		'l': {children: nil, actions: []actions.Action{actions.CursorRight}},

		'o': {children: nil, actions: []actions.Action{actions.InsertNewLineBelow, actions.Mode_Insert}},
		'O': {children: nil, actions: []actions.Action{actions.InsertNewLineAbove, actions.Mode_Insert}},

		'd': {actions: nil,
			children: map[keyboard.Key]*BindingNode{
				'd': {children: nil, actions: []actions.Action{actions.DeleteLine}},
			},
		},

		':': {children: nil, actions: []actions.Action{actions.Mode_Command}},
	},
}

var Insert = &BindingNode{
	actions: []actions.Action{actions.None},
	children: map[keyboard.Key]*BindingNode{
		keyboard.CtrlC: {children: nil, actions: []actions.Action{actions.Exit}},
		keyboard.ESC:   {children: nil, actions: []actions.Action{actions.Mode_Normal}},
		keyboard.ENTER: {children: nil, actions: []actions.Action{actions.InsertNewLine}},

		keyboard.ARROW_UP:    {children: nil, actions: []actions.Action{actions.CursorUp}},
		keyboard.ARROW_DOWN:  {children: nil, actions: []actions.Action{actions.CursorDown}},
		keyboard.ARROW_LEFT:  {children: nil, actions: []actions.Action{actions.CursorLeft}},
		keyboard.ARROW_RIGHT: {children: nil, actions: []actions.Action{actions.CursorRight}},

		keyboard.BACKSPACE:   {children: nil, actions: []actions.Action{actions.DeleteChar}},
		keyboard.BACKSPACE_2: {children: nil, actions: []actions.Action{actions.DeleteChar}},
		keyboard.DELETE:      {children: nil, actions: []actions.Action{actions.RemoveChar}},
		keyboard.TAB:         {children: nil, actions: []actions.Action{actions.InsertChar}},
	},
}

var Command = &BindingNode{
	actions: []actions.Action{},
	children: map[keyboard.Key]*BindingNode{
		keyboard.ESC:         {children: nil, actions: []actions.Action{actions.Mode_Normal, actions.ClearCommand}},
		keyboard.ENTER:       {children: nil, actions: []actions.Action{actions.RunCommand, actions.Submit}},
		keyboard.BACKSPACE:   {children: nil, actions: []actions.Action{actions.DeleteCommandChar}},
		keyboard.BACKSPACE_2: {children: nil, actions: []actions.Action{actions.DeleteCommandChar}},
	},
}
