//go:build debug
// +build debug

package logger

import (
	"log/slog"
	"os"
)

//func newLogger() *slog.Logger {
//	f, err := os.OpenFile("jte.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
//	if err != nil {
//		panic(err)
//	}
//	handler := slog.NewTextHandler(f, &slog.HandlerOptions{Level: slog.LevelDebug})
//	l := slog.New(handler)
//	return l
//}

func newLogger() *slog.Logger {
	f, err := os.OpenFile("jte.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	handler := slog.NewTextHandler(f, &slog.HandlerOptions{
		Level: slog.LevelDebug,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				// Remove the time attribute
				return slog.Attr{}
			}
			return a
		},
	})

	l := slog.New(handler)
	return l
}
