package commmand

import (
	"log/slog"
)

type CommandParser struct {
	logger *slog.Logger
}

func NewCommandParser(l *slog.Logger) *CommandParser {
	return &CommandParser{
		logger: l.WithGroup("command-parser"),
	}
}
// true if the command exists
func (cp *CommandParser) Parse(command string) (Command, bool) {
	c, exists := CommandMap[command]
	return c, exists
}
