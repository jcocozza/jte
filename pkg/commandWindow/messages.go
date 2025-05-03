package commandwindow

import (
	"fmt"
	"log/slog"
)

// a first in, first out
type MessageStack struct {
	logger *slog.Logger
	messages []string
}

func NewMessageStack(l *slog.Logger) *MessageStack {
	return &MessageStack{
		logger: l.WithGroup("message-queue"),
		messages: []string{},
	}
}

func (q *MessageStack) Size() int {
	return len(q.messages)
}

func (q *MessageStack) Push(msg string) {
	q.logger.Debug("push", slog.String("msg", msg))
	q.messages = append(q.messages, msg)
}

func (q *MessageStack) Pop() (string, error) {
	if len(q.messages) > 0 {
		idx := len(q.messages) - 1
		msg := q.messages[idx]
		q.logger.Debug("pop", slog.String("message", msg))
		q.messages = q.messages[:idx]
		return msg, nil
	}
	return "", fmt.Errorf("nothing to pop")
}

func (q *MessageStack) Clear() {
	q.messages = []string{}
}
