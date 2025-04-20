package bindings

import (
	"github.com/jcocozza/jte/pkg/actions"
	"github.com/jcocozza/jte/pkg/keyboard"
)

var Normal = &BindingNode{
	action: actions.None,
	children: map[keyboard.Key]*BindingNode{
		keyboard.CtrlC: {children: nil, action: actions.Exit},
		'w':            {children: nil, action: actions.None},
		'h':            {children: nil, action: actions.CursorLeft},
		'j':            {children: nil, action: actions.CursorDown},
		'k':            {children: nil, action: actions.CursorUp},
		'l':            {children: nil, action: actions.CursorRight},
	},
}
