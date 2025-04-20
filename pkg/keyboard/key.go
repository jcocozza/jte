package keyboard

type Key rune

func (k *Key) String() string {
	if s, ok := specialKeys[*k]; ok {
		return s
	}
	return string(*k)
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
