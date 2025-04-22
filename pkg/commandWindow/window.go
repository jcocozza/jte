package commandwindow

import (
	"log/slog"
	"strings"
)

type CommandWindow struct {
	mq      *MessageStack
	currcmd string
	ShowAll bool
	active  bool

	x int

	logger *slog.Logger
}

func NewCommandWindow(l *slog.Logger) *CommandWindow {
	return &CommandWindow{
		mq:      NewMessageStack(l),
		ShowAll: false,
		x:       0,

		logger: l.WithGroup("command-window"),
	}
}

func (c *CommandWindow) X() int {
	return c.x
}

func (c *CommandWindow) Reset() {
	c.logger.Debug("reset command window")
	c.mq.Clear()
	c.currcmd = ""
	c.x = 0
	c.active = false
}

func (c *CommandWindow) Activate() {
	c.logger.Debug("activate command window")
	c.active = true
}

func (c *CommandWindow) Active() bool {
	return c.active
}

func (c *CommandWindow) CmdBuf() string {
	return c.currcmd
}

// use the command registry to get the command
func (c *CommandWindow) ParseCommand() (Command, []string, error) {
	args := strings.Split(c.currcmd, " ")
	cmd := args[0]
	cmdArgs := args[1:]
	command, err := GetCommand(cmd)
	if err != nil {
		c.currcmd = ""
		return command, nil, err
	}
	return command, cmdArgs, nil
}

func (c *CommandWindow) Size() int {
	return c.mq.Size()
}

func (c *CommandWindow) AppendCharToCommand(char byte) {
	c.logger.Debug("append char", slog.String("char", string(char)))
	c.currcmd += string(char)
	c.x++
}

func (c *CommandWindow) RemoveCharFromCommand() {
	c.logger.Debug("remove char")
	if len(c.currcmd) > 0 {
		c.currcmd = c.currcmd[:len(c.currcmd)-1]
		c.x--
	}
}

func (c *CommandWindow) Push(content string) {
	c.logger.Debug("push", slog.String("content", content))
	c.mq.Push(content)
}

func (c *CommandWindow) PushMany(contents []string) {
	for _, msg := range contents {
		c.logger.Debug("push many", slog.String("content", msg))
		c.mq.Push(msg)
	}
}

// return the next thing to print
func (c *CommandWindow) Next() string {
	msg, err := c.mq.Pop()
	if err != nil {
		return ""
	}
	return msg
}

// return all content in the window
func (c *CommandWindow) Dump() []string {
	lst := []string{}
	for {
		msg, err := c.mq.Pop()
		if err != nil {
			break
		}
		lst = append(lst, msg)
	}
	return lst
}
