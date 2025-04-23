//go:build !release

package logger

import (
	"log/slog"
	"os"
)

var Logger *slog.Logger

func Init() (*os.File, error) {
	f, err := os.OpenFile("jte.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	handler := slog.NewTextHandler(f, &slog.HandlerOptions{Level: slog.LevelDebug})
	Logger = slog.New(handler)
	return f, nil
}

//func NewLogger() (*slog.Logger, *os.File, error) {
//	f, err := os.OpenFile("jte.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
//	if err != nil {
//		return nil, nil, err
//	}
//	handler := slog.NewTextHandler(f, &slog.HandlerOptions{Level: slog.LevelDebug})
//	logger := slog.New(handler)
//	return logger, f, nil
//}
