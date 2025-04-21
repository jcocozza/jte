package keyboard

import (
	"fmt"
	"log/slog"
)

type KeyQueue struct {
	logger *slog.Logger
	keys []Key
}

func NewKeyQueue(l *slog.Logger) *KeyQueue {
	return &KeyQueue{
		logger: l.WithGroup("key-queue"),
		keys: []Key{},
	}
}

func (q *KeyQueue) String() string {
	return string(q.keys)
}

func (q *KeyQueue) Enqueue(key Key) {
	q.logger.Debug("enqueue", slog.String("key", key.String()))
	q.keys = append(q.keys, key)
}

// pop the first element off the queue
func (q *KeyQueue) Dequeue() (Key, error) {
	if len(q.keys) > 0 {
		k := q.keys[0]
		q.logger.Debug("dequeue", slog.String("key", k.String()))
		q.keys = q.keys[1:]
		return k, nil
	}
	return Key(-1), fmt.Errorf("nothing to dequeue")
}
