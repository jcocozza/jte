//go:build !debug
// +build !debug

package logger

import "log/slog"

func newLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(nil, nil))
}
