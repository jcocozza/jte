package mode

import (
	"log/slog"
)

type mode string

const (
	ModeNavigation mode = "nav"
	ModeInsert     mode = "reg"
	ModeCommand    mode = "com"
)

type ModeManager struct {
	state mode
	logger *slog.Logger
}

func NewModeManager(l *slog.Logger) *ModeManager {
	return &ModeManager{state: ModeNavigation, logger: l.WithGroup("mode manager")}
}

func (m *ModeManager) Mode() mode {
	return m.state
}

func (m *ModeManager) SetMode(newMode mode) {
	m.logger.Info("set mode", slog.String("new mode", string(newMode)))
	m.state = newMode
}
