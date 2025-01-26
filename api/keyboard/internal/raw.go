// this handles raw presses from the keyboard
//
// it is mostly responsible for dealing with Special sequences like Ctrl+A, Arrow Keys, etc
package internal

import (
	"errors"
	"log/slog"
	"os"
	"unicode/utf8"
)

type SpecialKey uint16

const (
	F1 SpecialKey = 0xFFFF - iota
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
	//key_min // see terminfo
)

const (
	Ctrl_TILDE      SpecialKey = 0x00
	Ctrl2           SpecialKey = 0x00
	Ctrl_SPACE      SpecialKey = 0x00
	CtrlA           SpecialKey = 0x01
	CtrlB           SpecialKey = 0x02
	CtrlC           SpecialKey = 0x03
	CtrlD           SpecialKey = 0x04
	CtrlE           SpecialKey = 0x05
	CtrlF           SpecialKey = 0x06
	CtrlG           SpecialKey = 0x07
	BACKSPACE       SpecialKey = 0x08
	CtrlH           SpecialKey = 0x08
	TAB             SpecialKey = 0x09
	CtrlI           SpecialKey = 0x09
	CtrlJ           SpecialKey = 0x0A
	CtrlK           SpecialKey = 0x0B
	CtrlL           SpecialKey = 0x0C
	ENTER           SpecialKey = 0x0D
	CtrlM           SpecialKey = 0x0D
	CtrlN           SpecialKey = 0x0E
	CtrlO           SpecialKey = 0x0F
	CtrlP           SpecialKey = 0x10
	CtrlQ           SpecialKey = 0x11
	CtrlR           SpecialKey = 0x12
	CtrlS           SpecialKey = 0x13
	CtrlT           SpecialKey = 0x14
	CtrlU           SpecialKey = 0x15
	CtrlV           SpecialKey = 0x16
	CtrlW           SpecialKey = 0x17
	CtrlX           SpecialKey = 0x18
	CtrlY           SpecialKey = 0x19
	CtrlZ           SpecialKey = 0x1A
	ESC             SpecialKey = 0x1B
	Ctrl_LSQBRACKET SpecialKey = 0x1B
	Ctrl3           SpecialKey = 0x1B
	Ctrl4           SpecialKey = 0x1C
	Ctrl_BACKSLASH  SpecialKey = 0x1C
	Ctrl5           SpecialKey = 0x1D
	Ctrl_RSQBRACKET SpecialKey = 0x1D
	Ctrl6           SpecialKey = 0x1E
	Ctrl7           SpecialKey = 0x1F
	Ctrl_SLASH      SpecialKey = 0x1F
	Ctrl_UNDERSCORE SpecialKey = 0x1F
	//SPACE           SpecialKey = 0x20
	BACKSPACE_2 SpecialKey = 0x7F
	Ctrl8       SpecialKey = 0x7F
)

var specialKeys map[SpecialKey]string = map[SpecialKey]string{
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

// represents a single keypress
//
// a key press can either be a Special Key or a unicode character
type Keypress struct {
	Key     SpecialKey
	Unicode rune
}

func (kp *Keypress) IsUnicode() bool {
	return kp.Unicode != 0
}

// TODO: I'm not convinced that this method will work perfectly
func (kp *Keypress) String() string {
	if kp.Unicode == 0 {
		if name, ok := specialKeys[kp.Key]; ok {
			return name
		}
	}
	return string(kp.Unicode)
}

type Keyboard struct {
	logger *slog.Logger
}

func NewKeyboard(l *slog.Logger) *Keyboard {
	return &Keyboard{
		logger: l.WithGroup("raw keyboard"),
	}
}

// assumes that buf[0] = '\033'
// TODO: this is stupid logic. it needs to be improved and the checks should be more explicit
// they should not rely on the size of nread
func parseEscape(buf []byte, nread int) Keypress {
	if nread == 1 {
		return Keypress{Key: ESC}
	}
	if nread == 3 {
		switch buf[2] {
		case 'A':
			return Keypress{Key: ARROW_UP}
		case 'B':
			return Keypress{Key: ARROW_DOWN}
		case 'C':
			return Keypress{Key: ARROW_RIGHT}
		case 'D':
			return Keypress{Key: ARROW_LEFT}
		case 'F':
			return Keypress{Key: END}
		case 'H':
			return Keypress{Key: HOME}
		}
	}
	if nread == 4 {
		switch buf[2] {
		case '1', '7':
			return Keypress{Key: HOME}
		case '3':
			return Keypress{Key: DELETE}
		case '4', '8':
			return Keypress{Key: END}
		case '5':
			return Keypress{Key: PAGE_UP}
		case '6':
			return  Keypress{Key: PAGE_DOWN}
		}
	}
	return Keypress{Key: ESC}
}

func (kb *Keyboard) GetKeypress() (Keypress, error) {
	var buf = make([]byte, 4)
	var nread int
	var err error
	// TODO: instead of blocking for a for loop, maybe we can just use a channel
	for {
		nread, err = os.Stdin.Read(buf[:])
		if err != nil {
			return Keypress{}, err
		}
		if nread > 0 {
			break
		}
	}

	var kp Keypress
	switch {
	case buf[0] == '\033':
		kp = parseEscape(buf, nread)
	case SpecialKey(buf[0]) <= Ctrl_UNDERSCORE || SpecialKey(buf[0]) == BACKSPACE_2:
		kp = Keypress{Key: SpecialKey(buf[0])}
	default:
		r, _ := utf8.DecodeRune(buf[:nread])
		if r == utf8.RuneError {
			return Keypress{}, errors.New("invalid rune")
		}
		kp = Keypress{Unicode: r}
	}
	kb.logger.Info("keypress", slog.String("raw input", string(buf)), slog.String("key", kp.String()))
	return kp, nil
}
