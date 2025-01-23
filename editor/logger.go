package editor

import (
	"log/slog"
	"os"
)

func logFile() (*os.File, error) {
	path := "logs/editor-log.log"
	return os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
}

// create a logger
//
// input the desired log level
func CreateLogger(level slog.Level) *slog.Logger {
	opts := &slog.HandlerOptions{
		AddSource: true,
		Level:     level,
	}
	f, err := logFile()
	if err != nil {
		panic(err)
	}
	logger := slog.New(slog.NewTextHandler(f, opts))
	return logger
}
