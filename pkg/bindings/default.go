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
		'h':            {children: nil, action: actions.CursorLeft},
		'j':            {children: nil, action: actions.CursorDown},
		'k':            {children: nil, action: actions.CursorUp},
		'l':            {children: nil, action: actions.CursorRight},
	},
}

var Insert = &BindingNode{
	action: actions.None,
	children: map[keyboard.Key]*BindingNode{
		keyboard.CtrlC: {children: nil, action: actions.Exit},
		keyboard.ESC:   {children: nil, action: actions.Mode_Normal},
	},
}
