package action

import (
	"github.com/jcocozza/jte/internal/keyboard"
	"github.com/jcocozza/jte/internal/mode"
)

// this contains the default key bindings

var NormalBindings = &BindingNode{
	Actions: nil,
	children: map[keyboard.Key]*BindingNode{
		'i':            {children: nil, Actions: []Action{SwitchMode{m: mode.Insert}}},
		keyboard.CtrlC: {children: nil, Actions: []Action{Exit{}}},

		//'o': {children: nil, Actions: []Action{SwitchMode{m: mode.Insert}, NewLineBelow{}}},
		//'O': {children: nil, Actions: []Action{SwitchMode{m: mode.Insert}, NewLineAbove{}}},

		'd': {Actions: nil,
			children: map[keyboard.Key]*BindingNode{
				//'d': {children: nil, Actions: []Action{DeleteLine{}}},
			},
		},
		's': {children: nil, Actions: []Action{SplitHorizontal{}}},
		'v': {children: nil, Actions: []Action{SplitVertical{}}},
		'q': {children: nil, Actions: []Action{SplitClose{}}},

		'k': {children: nil, Actions: []Action{CursorUp{}}},
		'j': {children: nil, Actions: []Action{CursorDown{}}},
		'h': {children: nil, Actions: []Action{CursorLeft{}}},
		'l': {children: nil, Actions: []Action{CursorRight{}}},
	},
}
