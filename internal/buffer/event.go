package buffer

import (
	"fmt"
	"log/slog"
)

type EventType int

const (
	Event_Insert EventType = iota
	Event_Delete
	Event_Replace
)

// an event is a "block" of changes grouped together
//
// not 100% sure how the grouping will take place
type Event struct {
	complete bool
	etype    EventType
	changes  []Change
}

type EventStack []Event

func (s *EventStack) Push(e Event) {
	*s = append(*s, e)
}

func (s *EventStack) Pop() (Event, error) {
	if len(*s) == 0 {
		return Event{}, fmt.Errorf("stack is empty")
	}
	res := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return res, nil
}

type EventManager struct {
	history EventStack
	redo    EventStack
	current *Event

	logger *slog.Logger
}

func NewEventManager(l *slog.Logger) *EventManager {
	return &EventManager{
		history: EventStack{},
		redo: EventStack{},
		current: nil,
		logger: l.WithGroup("event-manager"),
	}
}

func (e *EventManager) Commit() {
	e.current.complete = true
	e.history.Push(*e.current)
	e.current = nil
}

func (e *EventManager) StartEvent(etype EventType) error {
	if e.current != nil {
		return fmt.Errorf("event already in progress")
	}
	e.current = &Event{
		complete: false,
		etype:    etype,
		changes:  []Change{},
	}
	return nil
}

func (e *EventManager) AddChange(c Change) {
	e.current.changes = append(e.current.changes, c)
}
