package editor

import (
	"bufio"
	"bytes"
	"fmt"
	"log/slog"
	"os"
	"unnamed/term"
)

type Editor struct {
	screenrows int
	screencols int

	abuf abuf
	rw   *term.RawMode

	c *cursor

	rows []*erow

	rowoffset int
	coloffset int

	logger *slog.Logger
}

func InitEditor() *Editor {
	rw, err := term.EnableRawMode()
	if err != nil {
		panic(err)
	}
	r, c, err := rw.WindowSize()
	if err != nil {
		panic(err)
	}
	e := &Editor{
		screenrows: r,
		screencols: c,
		abuf:       abuf{},
		rw:         rw,
		c:          &cursor{},
		rows:       []*erow{},
		logger:     CreateLogger(slog.LevelDebug),
	}
	e.logger.Info("editor initialized set up")
	return e
}

func (e *Editor) appendRow(row []byte) {
	er := &erow{
		chars:  row,
		render: []byte{},
	}
	er.Render()
	e.rows = append(e.rows, er)
}

func (e *Editor) Open(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Bytes()
		line = bytes.TrimRight(line, "\r\n")
		e.appendRow(line)
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}
	return nil
}

func (e *Editor) Exit(msg string) {
	if e.rw != nil {
		err := e.rw.Restore()
		if err != nil {
			e.logger.Error("failed to restore terminal %s; %s", msg, err.Error())
		}
	}
	e.logger.Info(msg)
	os.Exit(0)
}

func (e *Editor) ExitErr(err error) {
	if e.rw != nil {
		errR := e.rw.Restore()
		if errR != nil {
			e.logger.Error("failed to restore terminal: %s", errR.Error(), errR)
		}
	}
	e.logger.Error("exit with error", err.Error(), err)
	fmt.Fprintln(os.Stderr, "error: ", err)
	os.Exit(1)
}

// Modified readKeypress to handle raw bytes
func (e *Editor) readKeypress() rune {
	var buf [1]byte
	for {
		nread, err := os.Stdin.Read(buf[:])
		if nread == 1 {
			break
		}
		if err != nil {
			e.ExitErr(err)
		}
	}
	if rune(buf[0]) != EscapeSequence {
		e.logger.Info("read keypress", slog.String("byte", string(buf[0])))
		return rune(buf[0])
	}
	// Handle escape sequence
	var seq [2]byte
	n, err := os.Stdin.Read(seq[:1])
	if err != nil || n != 1 {
		return EscapeSequence
	}
	n, err = os.Stdin.Read(seq[1:2])
	if err != nil || n != 1 {
		return EscapeSequence
	}
	if seq[0] == '[' {
		if seq[1] >= '0' && seq[1] <= '9' {
			// Handle numeric escape sequence
			var third [1]byte
			n, err := os.Stdin.Read(third[:])
			if err != nil || n != 1 {
				return EscapeSequence
			}
			if third[0] == '~' {
				switch seq[1] {
				case '1', '7':
					return HOME
				case '3':
					return DELETE
				case '4', '8':
					return END
				case '5':
					return PAGE_UP
				case '6':
					return PAGE_DOWN
				}
			}
		}
		switch seq[1] {
		case 'A':
			return ARROW_UP
		case 'B':
			return ARROW_DOWN
		case 'C':
			return ARROW_RIGHT
		case 'D':
			return ARROW_LEFT
		case 'F':
			return END
		case 'H':
			return HOME
		}
	} else if seq[0] == 'O' {
		switch seq[1] {
		case 'H':
			return HOME
		case 'F':
			return END
		}
	}
	return EscapeSequence
}

// Modified ProcessKeypress to use the result directly
func (e *Editor) ProcessKeypress() {
	key := e.readKeypress()
	e.logger.Info("key read", slog.String("key", string(key)))

	if len(e.rows) == 0 && key != CtrlQ { // don't move when we have an empty file
		return
	}

	switch key {
	case CtrlQ:
		e.Exit("quit")
	case ARROW_UP:
		if e.c.Y > 0 {
			e.c.Y--
		}
	case ARROW_DOWN:
		if e.c.Y < len(e.rows)-1 {
			e.c.Y++
		}
	case ARROW_LEFT:
		if e.c.X > 0 {
			e.c.X--
		}
	case ARROW_RIGHT:
		if e.c.Y < len(e.rows) && e.c.X < len(e.rows[e.c.Y].render) {
			e.c.X++
		}
	case HOME:
		e.c.X = 0
	case END:
		if e.c.Y < len(e.rows) {
			e.c.X = len(e.rows[e.c.Y].render)
		}
	case PAGE_UP:
		e.c.Y = e.rowoffset
	case PAGE_DOWN:
		e.c.Y = e.rowoffset + e.screenrows - 1
		if e.c.Y > len(e.rows) {
			e.c.Y = len(e.rows)
		}
	default:
		printKey(key)
	}
	var newrow *erow
	if e.c.Y >= len(e.rows) {
		newrow = nil
	} else {
		newrow = e.rows[e.c.Y]
	}
	if e.c.X > len(newrow.render) {
		e.c.X = len(newrow.render)
	}
}

var welcome []byte = []byte("editor -- version: v0.0.1")

func (e *Editor) drawWelcome(rowNum int) {
	if rowNum == e.screenrows/3 {
		e.abuf.Append([]byte("~"))
		padding := (e.screencols - len(welcome)) / 2
		e.abuf.Append(bytes.Repeat([]byte(" "), padding))
		e.abuf.Append(welcome)
	} else {
		e.abuf.Append([]byte("~"))
	}

}

func (e *Editor) drawFile(rowNum int) {
	//e.abuf.Append([]byte("\r"))
	filerow := rowNum + e.rowoffset
	if filerow >= len(e.rows) {
		e.abuf.Append([]byte("~"))
	} else {
		e.abuf.Append(e.rows[filerow].render)
	}
}

func (e *Editor) draw() {
	for i := 0; i < e.screenrows; i++ {
		if len(e.rows) == 0 {
			e.drawWelcome(i)
		} else {
			e.drawFile(i)
		}
		e.abuf.Append([]byte("\x1b[K"))
		if i < e.screenrows-1 {
			e.abuf.Append([]byte("\r\n"))
		}
	}
}

func (e *Editor) drawCursor() {
	s := fmt.Sprintf("\x1b[%d;%dH", (e.c.Y-e.rowoffset)+1, (e.c.X-e.coloffset)+1)
	e.abuf.Append([]byte(s))
}

func (e *Editor) scroll() {
	if e.c.Y < e.rowoffset {
		e.rowoffset = e.c.Y
	}
	if e.c.Y >= e.rowoffset+e.screenrows {
		e.rowoffset = e.c.Y - e.screenrows + 1
	}
	if e.c.X < e.coloffset {
		e.coloffset = e.c.X
	}
	if e.c.X >= e.coloffset+e.screencols {
		e.coloffset = e.c.X - e.screencols + 1
	}

}

func (e *Editor) Refresh() {
	e.scroll()
	e.abuf.Append([]byte("\x1b[?25l"))
	e.abuf.Append([]byte("\x1b[H"))
	e.draw()
	e.drawCursor()
	e.abuf.Append([]byte("\x1b[?25h"))
	e.abuf.Flush()
}

func (e *Editor) Debug() {
	for _, row := range e.rows {
		fmt.Println(string(row.render))
	}
}
