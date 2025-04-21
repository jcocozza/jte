package commandwindow

import (
	"fmt"
	"log/slog"
)

type MessageQueue struct {
	logger *slog.Logger
	messages []string
}

func NewMessageQueue(l *slog.Logger) *MessageQueue {
	return &MessageQueue{
		logger: l.WithGroup("message-queue"),
		messages: []string{},
	}
}

func (q *MessageQueue) Size() int {
	return len(q.messages)
}

func (q *MessageQueue) Enqueue(msg string) {
	q.logger.Debug("enqueue", slog.String("msg", msg))
	q.messages = append(q.messages, msg)
}

func (q *MessageQueue) Dequeue() (string, error) {
	if len(q.messages) > 0 {
		msg := q.messages[0]
		q.logger.Debug("dequeue", slog.String("message", msg))
		q.messages = q.messages[1:]
		return msg, nil
	}
	return "", fmt.Errorf("nothing to dequeue")
}
