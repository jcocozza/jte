package command

import (
	"fmt"
	"log/slog"

	"github.com/jcocozza/jte/api/buffer"
	"github.com/jcocozza/jte/api/keyboard"
)

const (
	Command_LS = "ls"
)

var prompt = []byte("> ")

type CommandWindow struct {
	output [][]byte
	prompt []byte

	inputBuf []keyboard.Key

	previous []string
	logger   *slog.Logger
}

func NewCommandWindow(l *slog.Logger) *CommandWindow {
	return &CommandWindow{
		logger: l,
		prompt: []byte{},
	}
}

func (c *CommandWindow) NumRows() int {
	return len(c.output) + 1 // +1 for the prompt
}

func (c *CommandWindow) Output() [][]byte {
	return c.output
}
func (c *CommandWindow) Prompt() []byte {
	t := c.prompt
	for _, c := range c.inputBuf {
		t = append(t, byte(c))
	}
	return t
}

func (c *CommandWindow) Activate() {
	c.prompt = prompt
}

func (c *CommandWindow) AddInput(key keyboard.Key) {
	c.inputBuf = append(c.inputBuf, key)
}

func (c *CommandWindow) Handle(bm *buffer.BufferManager) {
	cmd := string(c.inputBuf)
	c.previous = append(c.previous, cmd)
	c.inputBuf = []keyboard.Key{}
	c.logger.Info("new command", slog.String("command", cmd))
	switch cmd {
	case Command_LS:
		bufLst := bm.List()
		for _, buf := range bufLst {
			c.output = append(c.output, []byte(fmt.Sprintf("%d: %s", buf.ID, buf.Name)))
		}
		c.prompt = []byte("any key to continue")
	default:
		c.prompt = []byte("invalid command (any key to continue)")
	}
}

func (c *CommandWindow) Clear() {
	c.output = [][]byte{}
	c.inputBuf = []keyboard.Key{}
	c.prompt = []byte{}
}
