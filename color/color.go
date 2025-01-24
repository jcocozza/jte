package color

import (
	"fmt"
	"unicode"
)

const (
	BLACK rune = iota + 30
	RED
	GREEN
	YELLOW
	BLUE
	MAGENTA
	CYAN
	WHITE
	UNKNOWN
	RESET
)

type Highlight int

const (
	HL_NORMAL Highlight = iota
	HL_NUMBER
	HL_MATCH
)

func ColorByte(b byte) Highlight {
	switch {
	case unicode.IsDigit(rune(b)):
		return HL_NUMBER
	default:
		return HL_NORMAL
	}
}

func MakeColor(hl rune) string {
	return fmt.Sprintf("\x1b[%dm", hl)
}

func SyntaxToColor(hl Highlight) string {
	switch hl {
	case HL_NUMBER:
		return MakeColor(RED)
	case HL_MATCH:
		return MakeColor(BLUE)
	case HL_NORMAL:
		fallthrough
	default:
		return MakeColor(WHITE)
	}
}
