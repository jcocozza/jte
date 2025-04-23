//go:build release

package logger

import (
	"log/slog"
	"os"
)

var Logger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))

func Init() (*os.File, error) {
	return nil, nil // no file to close
}
