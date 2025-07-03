package keyboard

import "unicode"

type Key rune

var specialKeys map[Key]string = map[Key]string{
	F1:          "F1",
	F2:          "F2",
	F3:          "F3",
	F4:          "F4",
	F5:          "F5",
	F6:          "F6",
	F7:          "F7",
	F8:          "F8",
	F9:          "F9",
	F10:         "F10",
	F11:         "F11",
	F12:         "F12",
	INSERT:      "INSERT",
	DELETE:      "DELETE",
	HOME:        "HOME",
	END:         "END",
	PAGE_UP:     "PAGE_UP",
	PAGE_DOWN:   "PAGE_DOWN",
	ARROW_UP:    "ARROW_UP",
	ARROW_DOWN:  "ARROW_DOWN",
	ARROW_LEFT:  "ARROW_LEFT",
	ARROW_RIGHT: "ARROW_RIGHT",
	Ctrl_TILDE:  "Ctrl+~",
	//Ctrl2:     "Ctrl+2",
	//CtrlSpace:     "Ctrl+SPACE",
	CtrlA:     "Ctrl+A",
	CtrlB:     "Ctrl+B",
	CtrlC:     "Ctrl+C",
	CtrlD:     "Ctrl+D",
	CtrlE:     "Ctrl+E",
	CtrlF:     "Ctrl+F",
	CtrlG:     "Ctrl+G",
	BACKSPACE: "BACKSPACE",
	//CtrlH:           "Ctrl+H",
	TAB: "TAB",
	//CtrlI:           "Ctrl+I",
	CtrlJ: "Ctrl+J",
	CtrlK: "Ctrl+K",
	CtrlL: "Ctrl+L",
	ENTER: "ENTER",
	//CtrlM:           "Ctrl+M",
	CtrlN: "Ctrl+N",
	CtrlO: "Ctrl+O",
	CtrlP: "Ctrl+P",
	CtrlQ: "Ctrl+Q",
	CtrlR: "Ctrl+R",
	CtrlS: "Ctrl+S",
	CtrlT: "Ctrl+T",
	CtrlU: "Ctrl+U",
	CtrlV: "Ctrl+V",
	CtrlW: "Ctrl+W",
	CtrlX: "Ctrl+X",
	CtrlY: "Ctrl+Y",
	CtrlZ: "Ctrl+Z",
	ESC:   "ESC",
	//Ctrl_LSQBRACKET: "Ctrl+[",
	//Ctrl3:           "Ctrl+3",
	Ctrl4: "Ctrl+4",
	//Ctrl_BACKSLASH:  "Ctrl+\\",
	Ctrl5: "Ctrl+5",
	//Ctrl_RSQBRACKET: "Ctrl+]",
	Ctrl6: "Ctrl+6",
	Ctrl7: "Ctrl+7",
	//Ctrl_SLASH:      "Ctrl+/",
	//Ctrl_UNDERSCORE: "Ctrl+_",
	//SPACE:       "SPACE",
	BACKSPACE_2: "BACKSPACE_2",
	//Ctrl8:           "Ctrl+8",
}

func (k *Key) String() string {
	if s, ok := specialKeys[*k]; ok {
		return s
	}
	return string(*k)
}

func (k Key) IsDigit() bool {
	return unicode.IsDigit(rune(k))
}

func (k Key) IsUnicode() bool {
	if k >= F1 && k <= Ctrl8 {
		return false
	}
	basicPlane := k >= 0x0000 && k <= 0xFFFF
	supplementaryPlane := k >= 0x100000 && k <= 0x10FFFF
	return basicPlane || supplementaryPlane
}

// using the unicode private use area for special keys
const (
	F1 Key = 0xE000 + iota
	F2
	F3
	F4
	F5
	F6
	F7
	F8
	F9
	F10
	F11
	F12
	INSERT
	DELETE
	HOME
	END
	PAGE_UP
	PAGE_DOWN
	ARROW_UP
	ARROW_DOWN
	ARROW_LEFT
	ARROW_RIGHT
	Ctrl_TILDE
	Ctrl2
	Ctrl_SPACE
	CtrlA
	CtrlB
	CtrlC
	CtrlD
	CtrlE
	CtrlF
	CtrlG
	BACKSPACE
	CtrlH
	TAB
	CtrlI
	CtrlJ
	CtrlK
	CtrlL
	ENTER
	CtrlM
	CtrlN
	CtrlO
	CtrlP
	CtrlQ
	CtrlR
	CtrlS
	CtrlT
	CtrlU
	CtrlV
	CtrlW
	CtrlX
	CtrlY
	CtrlZ
	ESC
	Ctrl_LSQBRACKET
	Ctrl3
	Ctrl4
	Ctrl_BACKSLASH
	Ctrl5
	Ctrl_RSQBRACKET
	Ctrl6
	Ctrl7
	Ctrl_SLASH
	Ctrl_UNDERSCORE
	BACKSPACE_2
	Ctrl8
)
