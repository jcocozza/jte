package messages

import "time"


var (
	MomentoMori Message = Message{"Momento Mori", time.Now(), time.Duration(3*time.Second)}
	Hello       Message = Message{"hello", time.Now(), -1}
	Goodbye     Message = Message{"good bye", time.Now(), -1}
	GoodDay     Message = Message{"good day", time.Now(), -1}
)

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
	return time.Now().After(m.Time.Add(m.Dur))
}

type MessageList []Message

func (q *MessageList) Push(msg Message) {
	*q = append(*q, msg)
}
