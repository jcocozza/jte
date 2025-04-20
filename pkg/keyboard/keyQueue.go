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
		logger: l,
		keys: nil,
	}
}

func (q *KeyQueue) Enqueue(key Key) {
	q.keys = append(q.keys, key)
}

// pop the first element off the queue
func (q *KeyQueue) Dequeue() (Key, error) {
	if len(q.keys) > 0 {
		k := q.keys[0]
		q.keys = q.keys[1:]
		return k, nil
	}
	return Key(-1), fmt.Errorf("nothing to dequeue")
}
