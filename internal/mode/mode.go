package mode

import (
	"fmt"
	"log/slog"
)

type Mode int

func (m *Mode) String() string {
	return modes[*m]
}

const (
	Normal Mode = iota
	Insert
	Command
)

var modes = [...]string{
	Normal:  "normal",
	Insert:  "insert",
	Command: "command",
}

type ModeMachine struct {
	current Mode
	logger  *slog.Logger
}

func NewModeMachine(l *slog.Logger) *ModeMachine {
	return &ModeMachine{logger: l.WithGroup("mode-machine")}
}

func (m *ModeMachine) Current() Mode {
	m.logger.Debug("current mode is " + modes[m.current])
	return m.current
}

func (mm *ModeMachine) SetMode(m Mode) {
	mm.logger.Debug(fmt.Sprintf("mode transition: %s => %s", modes[mm.current], modes[m]))
	mm.current = m
}
