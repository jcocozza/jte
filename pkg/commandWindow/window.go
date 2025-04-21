package commandwindow

import "log/slog"

type CommandWindow struct {
	mq      *MessageQueue
	ShowAll bool
}

func NewCommandWindow(l *slog.Logger) *CommandWindow {
	return &CommandWindow{
		mq:      NewMessageQueue(l),
		ShowAll: false,
	}
}

func (c *CommandWindow) Size() int {
	return c.mq.Size()
}

func (c *CommandWindow) Push(content string) {
	c.mq.Enqueue(content)
}

func (c *CommandWindow) PushMany(contents []string) {
	for _, msg := range contents {
		c.mq.Enqueue(msg)
	}
}

// return the next thing to print
func (c *CommandWindow) Next() string {
	msg, err := c.mq.Dequeue()
	if err != nil {
		return ""
	}
	return msg
}

// return all content in the window
func (c *CommandWindow) Dump() []string {
	lst := []string{}
	for {
		msg, err := c.mq.Dequeue()
		if err != nil {
			break
		}
		lst = append(lst, msg)
	}
	return lst
}
