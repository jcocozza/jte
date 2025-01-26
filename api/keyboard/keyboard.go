// an api that maps keyboard input into a unified data type (rune)
// exports non-unicode key presses as constants
package keyboard

import (
	"errors"
	"log/slog"

	"github.com/jcocozza/jte/api/keyboard/internal"
)

var ErrInvalidKey error = errors.New("invalid key")


type Keyboard struct {
	raw    internal.Keyboard
	logger *slog.Logger
}

func NewKeyboard(l *slog.Logger) *Keyboard {
	return &Keyboard{
		raw: *internal.NewKeyboard(l),
		logger: l.WithGroup("keyboard"),
	}
}

func (k *Keyboard) GetKeypress() (Key, error) {
	key, err := k.handleRawInput()
	if err != nil {
		return -1, err
	}
	k.logger.Info("keypress", slog.String("key", key.String()))
	return key, nil
}

func (k *Keyboard) handleRawInput() (Key, error) {
	kp, err := k.raw.GetKeypress()
	if err != nil {
		return -1, err
	}
	// this should handle most cases...most is unicode represented
	if kp.IsUnicode() {
		return Key(kp.Unicode), nil
	}
	switch kp.Key {
	case internal.F1:
		return F1, nil
	case internal.F2:
		return F2, nil
	case internal.F3:
		return F3, nil
	case internal.F4:
		return F4, nil
	case internal.F5:
		return F5, nil
	case internal.F6:
		return F6, nil
	case internal.F7:
		return F7, nil
	case internal.F8:
		return F8, nil
	case internal.F9:
		return F9, nil
	case internal.F10:
		return F10, nil
	case internal.F11:
		return F11, nil
	case internal.F12:
		return F12, nil
	case internal.INSERT:
		return INSERT, nil
	case internal.DELETE:
		return DELETE, nil
	case internal.HOME:
		return HOME, nil
	case internal.END:
		return END, nil
	case internal.PAGE_UP:
		return PAGE_UP, nil
	case internal.PAGE_DOWN:
		return PAGE_DOWN, nil
	case internal.ARROW_UP:
		return ARROW_UP, nil
	case internal.ARROW_DOWN:
		return ARROW_DOWN, nil
	case internal.ARROW_LEFT:
		return ARROW_LEFT, nil
	case internal.ARROW_RIGHT:
		return ARROW_RIGHT, nil
	case internal.Ctrl_TILDE:
		return Ctrl_TILDE, nil
	//case internal.Ctrl2:
	//return Ctrl2 , nil
	//case internal.Ctrl_SPACE:
	//return Ctrl_SPACE , nil
	case internal.CtrlA:
		return CtrlA, nil
	case internal.CtrlB:
		return CtrlB, nil
	case internal.CtrlC:
		return CtrlC, nil
	case internal.CtrlD:
		return CtrlD, nil
	case internal.CtrlE:
		return CtrlE, nil
	case internal.CtrlF:
		return CtrlF, nil
	case internal.CtrlG:
		return CtrlG, nil
	case internal.BACKSPACE:
		return BACKSPACE, nil
	//case internal.CtrlH:
	//return CtrlH , nil
	case internal.TAB:
		return TAB, nil
	//case internal.CtrlI:
	//return CtrlI , nil
	case internal.CtrlJ:
		return CtrlJ, nil
	case internal.CtrlK:
		return CtrlK, nil
	case internal.CtrlL:
		return CtrlL, nil
	case internal.ENTER:
		return ENTER, nil
	//case internal.CtrlM:
	//return CtrlM , nil
	case internal.CtrlN:
		return CtrlN, nil
	case internal.CtrlO:
		return CtrlO, nil
	case internal.CtrlP:
		return CtrlP, nil
	case internal.CtrlQ:
		return CtrlQ, nil
	case internal.CtrlR:
		return CtrlR, nil
	case internal.CtrlS:
		return CtrlS, nil
	case internal.CtrlT:
		return CtrlT, nil
	case internal.CtrlU:
		return CtrlU, nil
	case internal.CtrlV:
		return CtrlV, nil
	case internal.CtrlW:
		return CtrlW, nil
	case internal.CtrlX:
		return CtrlX, nil
	case internal.CtrlY:
		return CtrlY, nil
	case internal.CtrlZ:
		return CtrlZ, nil
	case internal.ESC:
		return ESC, nil
	//case internal.Ctrl_LSQBRACKET:
	//return Ctrl_LSQBRACKET , nil
	//case internal.Ctrl3:
	//return Ctrl3 , nil
	case internal.Ctrl4:
		return Ctrl4, nil
	//case internal.Ctrl_BACKSLASH:
	//return Ctrl_BACKSLASH , nil
	case internal.Ctrl5:
		return Ctrl5, nil
	//case internal.Ctrl_RSQBRACKET:
	//return Ctrl_RSQBRACKET , nil
	case internal.Ctrl6:
		return Ctrl6, nil
	case internal.Ctrl7:
		return Ctrl7, nil
	//case internal.Ctrl_SLASH:
	//return Ctrl_SLASH , nil
	//case internal.Ctrl_UNDERSCORE:
	//return Ctrl_UNDERSCORE , nil
	case internal.BACKSPACE_2:
		return BACKSPACE_2, nil
		//case internal.Ctrl8:
		//return Ctrl8 , nil
	}
	return -1, ErrInvalidKey
}
