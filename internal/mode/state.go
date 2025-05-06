package mode

import (
	"fmt"
	"log/slog"
)

type Mode string

const (
	Insert  Mode = "insert"
	Normal  Mode = "normal"
	Command Mode = "command"
)

type StateMachine struct {
	current Mode
	modes   map[Mode]struct{}
	logger  *slog.Logger
}

// defaults to Normal mode as current
func NewStateMachine(l *slog.Logger) *StateMachine {
	return &StateMachine{
		current: Normal,
		modes: map[Mode]struct{}{
			Insert:  {},
			Normal:  {},
			Command: {},
		},
		logger: l.WithGroup("state-machine"),
	}
}

func (s *StateMachine) SetMode(name Mode) {
	if _, ok := s.modes[name]; ok {
		s.logger.Debug("set mode", slog.String("mode", string(name)))
		s.current = name
		return
	}
	panic(fmt.Sprintf("unexpected mode: %s", name))
}

func (s *StateMachine) Current() Mode { return s.current }
