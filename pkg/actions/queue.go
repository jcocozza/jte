package actions

import (
	"fmt"
	"log/slog"
)

// a first in, first out queue
type ActionQueue struct {
	logger  *slog.Logger
	actions []Action
}

func NewActionQueue(l *slog.Logger) *ActionQueue {
	return &ActionQueue{
		logger: l.WithGroup("action-queue"),
		actions: []Action{},
	}
}

// add action to the queue
func (q *ActionQueue) Enqueue(action Action) {
	q.logger.Debug("enqueue", slog.String("action", ActionNames[action]))
	q.actions = append(q.actions, action)
}

// pop the next action
func (q *ActionQueue) Dequeue() (Action, error){
	if len(q.actions) > 0 {
		action := q.actions[0]
		q.logger.Debug("dequeue", slog.String("action", ActionNames[action]))
		q.actions = q.actions[1:]
		return action, nil
	}
	return None, fmt.Errorf("nothing to dequeue")
}

// process all actions in the queue
//func (q *ActionQueue) Process() {
//	for _, action := range q.actions {
//		q.logger.Debug("applying action", slog.String("action", action.Name()))
//		action.Apply()
//	}
//	// clear actions
//	q.actions = nil
//}
