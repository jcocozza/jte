package logger

import "log/slog"

func NewLogger() *slog.Logger {
	return newLogger()
}
