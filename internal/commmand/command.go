package commmand

import (
	"fmt"
	"log/slog"

	"github.com/jcocozza/jte/internal/keyboard"
)

type Command int

const (
	Empty Command = iota + 1
	Quit
	List
)

// map command string to command
var CommandMap = map[string]Command{
	"q":  Quit,
	"ls": List,
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
func (w *CommandWindow) GetCommand() (Command, error) {
	cmd := keyboard.Collapse(w.Input)
	command, _ := w.p.Parse(cmd)
	switch command {
	case Empty:
		w.ClearInput()
		w.ClearOutput()
		w.Hide()
	case Quit:
		w.Output = append(w.Output, "quitting...")
		w.Show()
		w.Lock()
		// TODO we need to actually quit
	case List:
		w.Output = append(w.Output, fmt.Sprintf("running command: %s", cmd))
		w.Output = append(w.Output, "foo")
		w.Output = append(w.Output, "foo")
		w.Output = append(w.Output, "foo")
		w.Output = append(w.Output, "foo")
		w.Output = append(w.Output, "<Esc> to continue.")
		w.Show()
		w.Lock()
	default:
		w.Output = append(w.Output, fmt.Sprintf("[ERROR] command %s does not exist", cmd))
		w.ClearInput()
		w.Show()
		return -1, fmt.Errorf("invalid command")
	}
	return command, nil
}
