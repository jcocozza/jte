package logger

import (
	"log/slog"
	"os"
)

func NewLogger() (*slog.Logger, *os.File, error) {
	f, err := os.OpenFile("jte.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, nil, err
	}
	handler := slog.NewTextHandler(f, nil)
	logger := slog.New(handler)
	return logger, f, nil
}
