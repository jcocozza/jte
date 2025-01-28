package command

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/jcocozza/jte/api/buffer"
	"github.com/jcocozza/jte/api/keyboard"
	"github.com/jcocozza/jte/api/messages"
	"github.com/jcocozza/jte/api/search"
)

const (
	Command_LS = "ls"
)

type CommandMode int
const (
	CommandInactive CommandMode = iota
	CommandBasic
	CommandSearch
)

var prompt = []byte("> ")
var searchPrompt = []byte("/")

type CommandWindow struct {
	output [][]byte
	prompt []byte

	inputBuf []keyboard.Key

	previous []string

	Mode CommandMode

	SearchResults *search.SearchResults

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

func (c *CommandWindow) ShrinkInput() {
	if len(c.inputBuf) == 0 {
		return
	}
	c.inputBuf = c.inputBuf[0:len(c.inputBuf)-1]
}

func (c *CommandWindow) Activate() {
	c.prompt = prompt
	c.Mode = CommandBasic
}

func (c *CommandWindow) ActivateSearch() {
	c.prompt = searchPrompt
	c.Mode = CommandSearch
}

func (c *CommandWindow) AddInput(key keyboard.Key) {
	c.inputBuf = append(c.inputBuf, key)
}

func (c *CommandWindow) SetMessage(msg messages.Message) {
	if c.Mode != CommandInactive {
		return
	}
	c.inputBuf = []keyboard.Key{}
	c.prompt = []byte(msg.Text)
	if msg.Dur != -1 {
		go func() { // I don't think this is a good way to do this, but for now, it will work
			time.Sleep(msg.Dur)
			c.SetMessage(messages.Message{})
		}()
	}
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
	c.Mode = CommandInactive
}

func (c *CommandWindow) SearchPattern() string {
	if c.Mode != CommandSearch {
		return ""
	}
	return string(c.inputBuf)
}

func (c *CommandWindow) HandleSearch(buf buffer.Buffer) {
	pattern := string(c.inputBuf)
	c.logger.Debug("searching", slog.String("pattern", pattern), slog.Int("in buf len", len(c.inputBuf)))
	c.SearchResults = search.SearchItr(pattern, buf)
}

func (c *CommandWindow) Clear() {
	c.output = [][]byte{}
	c.inputBuf = []keyboard.Key{}
	c.prompt = []byte{}
}
