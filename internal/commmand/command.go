package commmand

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/jcocozza/jte/internal/keyboard"
)

type Command int

const (
	Empty Command = iota + 1
	Quit
	List
	Edit
)

// map command string to command
var CommandMap = map[string]Command{
	"q":    Quit,
	"quit": Quit,

	"ls":   List,
	"list": List,

	"e":    Edit,
	"edit": Edit,
}

type CommandWindow struct {
	logger *slog.Logger
	Input  []keyboard.Key
	Output []string
	p      *CommandParser

	ShowOutput bool
	locked     bool
}

func NewCommandWindow(l *slog.Logger) *CommandWindow {
	return &CommandWindow{
		logger: l.WithGroup("command-window"),
		Input:  []keyboard.Key{},
		Output: []string{},
		p:      NewCommandParser(l),
	}
}

func (w *CommandWindow) Push(msg string) {
	w.Output = append(w.Output, msg)
	w.Show()
}

func (w *CommandWindow) PushErr(err error) {
	msg := "[ERROR] " + err.Error()
	w.Output = []string{msg}
	w.Show()
}

func (w *CommandWindow) ClearAndPush(contents []string) {
	w.Output = contents
	w.Show()
}

func (w *CommandWindow) Locked() bool {
	return w.locked
}

func (w *CommandWindow) Hide() {
	w.ShowOutput = false
}
func (w *CommandWindow) Show() {
	w.ShowOutput = true
}

func (w *CommandWindow) ClearInput() {
	w.Input = []keyboard.Key{}
}
func (w *CommandWindow) ClearOutput() {
	w.Output = []string{}
}

func (w *CommandWindow) AddInput(k keyboard.Key) {
	if !w.locked {
		w.Input = append(w.Input, k)
	}
}

func (w *CommandWindow) Lock() {
	w.locked = true
}

func (w *CommandWindow) Unlock() {
	w.locked = false
}

// TODO: this is not a good way to do commands
func (w *CommandWindow) GetCommand() (Command, []string, error) {
	input := keyboard.Collapse(w.Input)
	if len(input) == 0 {
		return -1, nil, fmt.Errorf("no command")
	}
	inputSplit := strings.Split(input, " ")
	cmd, args := inputSplit[0], []string{}
	if len(inputSplit) > 1 {
		args = inputSplit[1:]
	}
	command, isreal := w.p.Parse(cmd)
	if isreal {
		return command, args, nil
	} else {
		return -1, args, fmt.Errorf("invalid command")
	}

	/* // keeping this for now because i need to figure out a better system entirely
	switch command {
	case Empty:
		w.ClearInput()
		w.ClearOutput()
		w.Hide()
	case Quit:
		w.Output = append(w.Output, "quitting...")
		w.Show()
		w.Lock()
	case List:
		w.Output = append(w.Output, fmt.Sprintf("running command: %s", cmd))
		w.Output = append(w.Output, "foo")
		w.Output = append(w.Output, "foo")
		w.Output = append(w.Output, "foo")
		w.Output = append(w.Output, "foo")
		w.Output = append(w.Output, "<Esc> to continue.")
		w.Show()
		w.Lock()
	case Edit:
	default:
		w.Output = append(w.Output, fmt.Sprintf("[ERROR] command %s does not exist", cmd))
		w.ClearInput()
		w.Show()
		return -1, args, fmt.Errorf("invalid command")
	}
	return command, args, nil
	*/
}
