package editor

import (
	"fmt"
)

const (
	Null           = 0x00
	CtrlA          = 0x01
	CtrlB          = 0x02
	CtrlC          = 0x03
	CtrlD          = 0x04
	CtrlE          = 0x05
	CtrlF          = 0x06
	CtrlG          = 0x07
	CtrlH          = 0x08
	CtrlI          = 0x09
	TAB            = CtrlI
	CtrlJ          = 0x0A
	CtrlK          = 0x0B
	CtrlL          = 0x0C
	CtrlM          = 0x0D
	CarriageReturn = CtrlM
	CtrlN          = 0x0E
	CtrlO          = 0x0F
	CtrlP          = 0x10
	CtrlQ          = 0x11

	EscapeSequence rune = '\x1b'

	PAGE_UP rune = iota + 1000
	PAGE_DOWN
	ARROW_UP
	ARROW_DOWN
	ARROW_LEFT
	ARROW_RIGHT

	HOME
	END

	DELETE
)

func controlCharacterName(c rune) string {
	switch c {
	case 0x00:
		return "NULL"
	case 0x01:
		return "Ctrl-A"
	case 0x02:
		return "Ctrl-B"
	case 0x03:
		return "Ctrl-C"
	case 0x04:
		return "Ctrl-D"
	case 0x05:
		return "Ctrl-E"
	case 0x06:
		return "Ctrl-F"
	case 0x07:
		return "Ctrl-G"
	case 0x08:
		return "Ctrl-H"
	case 0x09:
		return "Tab"
	case 0x0A:
		return "Ctrl-J (Line Feed)"
	case 0x0B:
		return "Ctrl-K"
	case 0x0C:
		return "Ctrl-L"
	case 0x0D:
		return "Ctrl-M (Carriage Return)"
	case 0x0E:
		return "Ctrl-N"
	case 0x0F:
		return "Ctrl-O"
	default:
		return fmt.Sprintf("Ctrl-%c", 'A'+(c-1))
	}
}

func printKey(key rune) {
	if key <= 0x1F || key == 0x7F { // Check for control characters
		fmt.Printf("Control character: %d (%s)\n", key, controlCharacterName(key))
	} else {
		fmt.Printf("Printable character: %d::%s\r\n", key, string(key))
	}
}
