package keyboard

import (
	"errors"
	"log/slog"
	"os"
	"unicode/utf8"
)

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
		logger: l,
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
	case SpecialKey(buf[0]) <= SPACE || SpecialKey(buf[0]) == BACKSPACE_2:
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
