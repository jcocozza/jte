package editor

import (
	//"fmt"
	"log/slog"
	"os"
	//"time"
)

func logFile() (*os.File, error) {
	//path := fmt.Sprintf("logs/editor-log.log", time.Now().Unix())
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
