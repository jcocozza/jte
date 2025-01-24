package messages

import "time"

type Message struct {
	Text string
	Time time.Time
	Dur  time.Duration
}

func (m *Message) NonEmpty() bool {
	return m.Text != "" && !m.Time.IsZero()
}

func (m *Message) Expired() bool {
	if m.Dur == -1 {
		return false
	}
	return m.Time.Add(m.Dur).After(time.Now())
}

type MessageList []Message

func (q *MessageList) Push(msg Message) {
	*q = append(*q, msg)
}
