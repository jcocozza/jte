package actions

import (
	"log/slog"
)

// a first in, first out queue
type ActionQueue struct {
	logger  *slog.Logger
	actions []Action
}

func NewActionQueue(l *slog.Logger) *ActionQueue {
	return &ActionQueue{
		logger: l,
		actions: []Action{},
	}
}

// add action to the queue
func (q *ActionQueue) Enqueue(action Action) {
	q.actions = append(q.actions, action)
}

// apply the latest next action
func (q *ActionQueue) Dequeue() {
	if len(q.actions) > 0 {
		q.actions[0].Apply()
	}
	q.actions = q.actions[1:]
}

// process all actions in the queue
func (q *ActionQueue) Process() {
	for _, action := range q.actions {
		q.logger.Debug("applying action", slog.String("action", action.Name()))
		action.Apply()
	}
	// clear actions
	q.actions = nil
}
