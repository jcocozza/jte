package buffer

import "fmt"

type EventType int

const (
	Event_Insert EventType = iota
	Event_Remove
	Event_Replace
)

// a change is the result of a single key press
//
// on delete, the "contents" is the deleted text
// on insert, the "contents" is the inserted text
// on replace, the "contents" is the replaced text (similar to delete)
type change struct {
	startCur Cursor
	endCur   Cursor
	contents []byte
}

// an event is a "block" of changes grouped together
//
// not 100% sure how the grouping will take place
type Event struct {
	etype   EventType
	changes []change
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
