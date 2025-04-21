package bindings

import (
	"github.com/jcocozza/jte/pkg/actions"
	"github.com/jcocozza/jte/pkg/keyboard"
)

var Normal = &BindingNode{
	action: actions.None,
	children: map[keyboard.Key]*BindingNode{
		keyboard.CtrlC: {children: nil, action: actions.Exit},
		'i':            {children: nil, action: actions.Mode_Insert},

		'h': {children: nil, action: actions.CursorLeft},
		'j': {children: nil, action: actions.CursorDown},
		'k': {children: nil, action: actions.CursorUp},
		'l': {children: nil, action: actions.CursorRight},

		'o': {children: nil, action: actions.InsertNewLineBelow},
		'O': {children: nil, action: actions.InsertNewLineAbove},
	},
}

var Insert = &BindingNode{
	action: actions.None,
	children: map[keyboard.Key]*BindingNode{
		keyboard.CtrlC: {children: nil, action: actions.Exit},
		keyboard.ESC:   {children: nil, action: actions.Mode_Normal},
		keyboard.ENTER: {children: nil, action: actions.InsertNewLine},

		keyboard.ARROW_UP:    {children: nil, action: actions.CursorUp},
		keyboard.ARROW_DOWN:  {children: nil, action: actions.CursorDown},
		keyboard.ARROW_LEFT:  {children: nil, action: actions.CursorLeft},
		keyboard.ARROW_RIGHT: {children: nil, action: actions.CursorRight},

		keyboard.BACKSPACE:   {children: nil, action: actions.DeleteChar},
		keyboard.BACKSPACE_2: {children: nil, action: actions.DeleteChar},
		keyboard.DELETE:      {children: nil, action: actions.RemoveChar},
		keyboard.TAB:         {children: nil, action: actions.InsertChar},
	},
}
