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
		':': 			{children: nil, Actions: []Action{CommandClearOutput{},SwitchMode{m: mode.Command}}},
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

		'k':                  {children: nil, Actions: []Action{CursorUp{}}},
		'j':                  {children: nil, Actions: []Action{CursorDown{}}},
		'h':                  {children: nil, Actions: []Action{CursorLeft{}}},
		'l':                  {children: nil, Actions: []Action{CursorRight{}}},
		keyboard.ARROW_UP:    {children: nil, Actions: []Action{CursorUp{}}},
		keyboard.ARROW_DOWN:  {children: nil, Actions: []Action{CursorDown{}}},
		keyboard.ARROW_LEFT:  {children: nil, Actions: []Action{CursorLeft{}}},
		keyboard.ARROW_RIGHT: {children: nil, Actions: []Action{CursorRight{}}},

		leader: {
			Actions: nil,
			children: map[keyboard.Key]*BindingNode{
				'k': {children: nil, Actions: []Action{PaneUp{}}},
				'j': {children: nil, Actions: []Action{PaneDown{}}},
				'h': {children: nil, Actions: []Action{PaneLeft{}}},
				'l': {children: nil, Actions: []Action{PaneRight{}}},
			},
		},
	},
}

var InsertBindings = &BindingNode{
	Actions: nil,
	children: map[keyboard.Key]*BindingNode{
		keyboard.ESC:         {children: nil, Actions: []Action{SwitchMode{m: mode.Normal}}},
		keyboard.CtrlC:       {children: nil, Actions: []Action{Exit{}}},
		keyboard.ARROW_UP:    {children: nil, Actions: []Action{CursorUp{}}},
		keyboard.ARROW_DOWN:  {children: nil, Actions: []Action{CursorDown{}}},
		keyboard.ARROW_LEFT:  {children: nil, Actions: []Action{CursorLeft{}}},
		keyboard.ARROW_RIGHT: {children: nil, Actions: []Action{CursorRight{}}},
	},
}

var CommandBindings = &BindingNode{
	Actions: nil,
	children: map[keyboard.Key]*BindingNode{
		keyboard.ESC: {children: nil, Actions: []Action{CommandClearOutput{}, CommandClearInput{}, SwitchMode{m: mode.Normal}}},
		keyboard.CtrlC: {children: nil, Actions: []Action{Exit{}}},
		keyboard.ENTER: {children: nil, Actions: []Action{CommandRun{}}},
	},
}
