package editor

import (
	"bufio"
	"bytes"
	"fmt"
	"log/slog"
	"os"
	"time"
	"unnamed/fileutil"
	"unnamed/term"
)

const msgTimeout = time.Duration(3 * time.Second)

type Editor struct {
	screenrows int
	screencols int

	abuf abuf
	rw   *term.RawMode
	c    *cursor
	rows []*erow

	rowoffset int
	coloffset int

	filename string
	// has the file been modified since being opened
	dirty bool

	// todo: this seems like a perfect usecase for go channels
	msg     string
	msgTime time.Time

	// if we are on the welcome page
	welcome bool

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
		screenrows: r - 2, // leave room for status bar and messages
		screencols: c,
		abuf:       abuf{},
		rw:         rw,
		c:          &cursor{},
		rows:       []*erow{},
		logger:     CreateLogger(slog.LevelDebug),
		welcome:    true,
	}
	e.logger.Info("editor initialized set up", slog.String("location", fmt.Sprintf("%d, %d", e.c.X, e.c.Y)))
	return e
}

func (e *Editor) status() string {
	var displayName string = e.filename
	if e.filename == "" {
		displayName = "[No Name]"
	}
	var displayDirty string = ""
	if e.dirty {
		displayDirty = "(modified)"
	}
	return fmt.Sprintf("ln:%d/%d - %s %s", e.c.Y, len(e.rows)-1, displayDirty, displayName)
}

func (e *Editor) SetMsg(msg string, timeout time.Duration) {
	e.msgTime = time.Now()
	e.msg = msg
	e.Refresh() // this is need to show the prompt
	go func() {
		if timeout == -1 {
			return
		}
		e.msg = ""
		e.msgTime = time.Time{}
	}()
}

// pass in a nil callback to just receive the user input back
// otherwise will call the callback
func (e *Editor) prompt(prompt string, callback func(input []byte, char rune)) []byte {
	e.SetMsg(prompt, -1)
	var buf []byte
	runCallback := func(key rune) {
		if callback != nil {
			callback(buf, key)
		}
	}
	for {
		key := e.readKeypress()
		switch {
		case key == DEL:
			fallthrough
		case key == CtrlH:
			fallthrough
		case key == BACKSPACE:
			if len(buf) != 0 {
				buf = buf[:len(buf)-1]
				msg := fmt.Sprintf("%s%s", prompt, string(buf))
				e.SetMsg(msg, -1)
			}
		case key == EscapeSequence:
			e.SetMsg("", -1)
			runCallback(key)
			return []byte{}
		case key == '\r':
			if len(buf) != 0 {
				e.SetMsg("", -1)
				runCallback(key)
				return buf
			}
		case key != CtrlC && key < 128:
			buf = append(buf, byte(key))
			msg := fmt.Sprintf("%s%s", prompt, string(buf))
			e.SetMsg(msg, -1)
		}
		runCallback(key)
	}
}

func (e *Editor) findCallback(query []byte, key rune) {
	if (key == '\r' || key == EscapeSequence) {
		return
	}
	for y, row := range e.rows {
		if bytes.Contains(row.render, query) {
			i := bytes.LastIndex(row.render, query)
			if i == -1 {
				return
			}
			e.c.Y = y
			e.c.X = i
			e.rowoffset = len(e.rows)
			return
		}
	}
}

func (e *Editor) find() {
	query := e.prompt("search (esc to cancel): ", e.findCallback)
	if len(query) == 0 {
		return
	}
}

func (e *Editor) insertRow(at int, row []byte) {
	if at < 0 || at > len(e.rows) {
		return
	}

	newRows := make([]*erow, len(e.rows)+1)
	copy(newRows[:at], e.rows[:at])
	newRows[at] = &erow{
		chars:  row,
		render: []byte(""),
	}
	newRows[at].Render()
	copy(newRows[at+1:], e.rows[at:])
	e.rows = newRows
}

func (e *Editor) deleteRow(at int) {
	if at < 0 || at >= len(e.rows) {
		return
	}
	e.rows = append(e.rows[:at], e.rows[at+1:]...)
}

func (e *Editor) clear() {
	e.welcome = false
	e.rows = []*erow{}
	e.abuf.Clear()
	e.c.reset()
	e.rowoffset = 0
	e.coloffset = 0
}

func (e *Editor) Open(filename string) error {
	file, err := fileutil.OpenOrCreateFile(filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	scanner := bufio.NewScanner(file)
	numScans := 0
	for scanner.Scan() {
		line := scanner.Bytes()
		line = bytes.TrimRight(line, "\r\n")
		e.insertRow(len(e.rows), line)
		numScans += 1
	}
	if numScans == 0 {
		e.insertRow(0, []byte{})
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}
	e.filename = filename
	e.dirty = false
	e.welcome = false
	return nil
}

func (e *Editor) save() error {
	if e.filename == "" {
		e.filename = string(e.prompt("save as (ESC to cancel): ", nil))
		if e.filename == "" {
			e.SetMsg("save aborted", msgTimeout)
			return nil
		}
	}
	buf := e.combineRows()
	n, err := fileutil.Save(e.filename, buf)
	if err != nil {
		return err
	}
	e.SetMsg(fmt.Sprintf("%d bytes written", n), msgTimeout)
	e.dirty = false
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

func (e *Editor) insertChar(c byte) {
	if e.c.Y == len(e.rows) {
		e.abuf.Append([]byte{})
		e.insertRow(len(e.rows), []byte{})
	}
	e.rows[e.c.Y].InsertChar(e.c.X, c)
	e.c.X++
	e.dirty = true
}

func (e *Editor) insertNewline() {
	if e.c.X == 0 {
		e.insertRow(e.c.Y, []byte{})
	} else {
		e.insertRow(e.c.Y+1, e.rows[e.c.Y].chars[e.c.X:])
		e.rows[e.c.Y].Trim(e.c.X)
	}
	e.c.Y++
	e.c.X = 0
}

func (e *Editor) deleteChar() {
	e.logger.Info("deleting char", slog.String("location", fmt.Sprintf("%d, %d", e.c.X, e.c.Y)))
	if e.c.Y == len(e.rows) {
		return
	}
	if e.c.X == 0 && e.c.Y == 0 {
		return
	}
	if e.c.X > 0 {
		e.rows[e.c.Y].DelChar(e.c.X - 1)
		e.c.X--
	} else {
		newX := len(e.rows[e.c.Y-1].chars)
		e.rows[e.c.Y-1].append(e.rows[e.c.Y].chars)
		e.deleteRow(e.c.Y)
		e.c.Y--
		e.c.X = newX
	}
	e.dirty = true
}

func (e *Editor) combineRows() []byte {
	var buf bytes.Buffer
	for _, row := range e.rows {
		buf.Write(row.chars)
		buf.WriteByte('\n')
	}
	return buf.Bytes()
}

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
	if rune(buf[0]) == DEL {
		return BACKSPACE
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
					return DEL
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
	e.logger.Info("key read before", slog.String("key", string(key)), slog.String("location", fmt.Sprintf("%d, %d", e.c.X, e.c.Y)))
	// don't move when we have an empty file
	// allow quitting and new file creation
	if len(e.rows) == 0 && key != CtrlQ && key != CtrlN {
		return
	}
	switch key {
	case '\r':
		e.insertNewline()
	case CtrlQ:
		if e.dirty {
			e.SetMsg("file has unsaved changes", msgTimeout)
			break
		}
		e.abuf.Append([]byte("\x1b[2J"))
		e.abuf.Append([]byte("\x1b[H"))
		e.Exit("quit")
	case CtrlS:
		err := e.save()
		if err == fileutil.ErrNoFilename {
			e.SetMsg("cannot save. no filename", msgTimeout)
		}
		if err != nil {
			e.ExitErr(err)
		}
		e.SetMsg("saved", msgTimeout)
	case CtrlN:
		if e.dirty {
			e.SetMsg("file has unsaved changes", msgTimeout)
			break
		}
		e.logger.Info("should open new file")
		e.clear()
		err := e.Open("")
		if err != nil {
			e.ExitErr(err)
		}
	case CtrlF:
		e.find()
	case CtrlP: // this is just for testing
		msg := e.prompt("prompt: ", nil)
		e.SetMsg(string(msg), msgTimeout)
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
	case BACKSPACE:
		e.logger.Info("backspace pressed", slog.String("location", fmt.Sprintf("%d, %d", e.c.X, e.c.Y)))
		if e.c.X > 0 || e.c.Y > 0 {
			e.deleteChar()
		}
	case DEL:
		e.logger.Info("delete pressed", slog.String("location", fmt.Sprintf("%d, %d", e.c.X, e.c.Y)))
		if e.c.Y < len(e.rows) && e.c.X < len(e.rows[e.c.Y].render) {
			e.c.X++
		}
		e.deleteChar()
	case CtrlL:
	case EscapeSequence:
		break
	default:
		e.insertChar(byte(key))
		//printKey(key)
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
	e.logger.Info("key read after", slog.String("key", string(key)), slog.Int("key int", int(key)), slog.String("location", fmt.Sprintf("%d, %d", e.c.X, e.c.Y)))
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

// draw status bar and any messages
func (e *Editor) drawStatusBar() {
	// status bar
	e.abuf.Append([]byte("\x1b[7m"))
	status := e.status()
	e.abuf.Append(bytes.Repeat([]byte(" "), e.screencols-len(status)))
	e.abuf.Append([]byte(status))
	e.abuf.Append([]byte("\x1b[m"))
	e.abuf.Append([]byte("\r\n"))

	// messages
	e.abuf.Append([]byte("\x1b[K"))
	var displayMsg string = e.msg
	if len(e.msg) > e.screencols {
		displayMsg = e.msg[0:e.screencols]
	}
	if e.msg != "" && !e.msgTime.IsZero() {
		e.abuf.Append([]byte(displayMsg))
	}
}

func (e *Editor) draw() {
	for i := 0; i < e.screenrows; i++ {
		if e.welcome {
			e.drawWelcome(i)
		} else {
			e.drawFile(i)
		}
		e.abuf.Append([]byte("\x1b[K"))
		e.abuf.Append([]byte("\r\n"))
	}
	e.drawStatusBar()
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
