package commandwindow

import "fmt"

type Command int

// steps for adding a command:
// 1. create a new const
// 2. add command to command registry (here)
// 3. add command to command function registry (pkg/editor/registry.go)

const (
	LS Command = iota
	ECHO
)

// map commands to a list of actions
var commandRegistry = map[string]Command{
	"ls": LS,
	"echo": ECHO,
}

// return an error if the command is not found
func GetCommand(cmd string) (Command, error) {
	if command, ok := commandRegistry[cmd]; ok {
		return command, nil
	}
	return -1, fmt.Errorf("command %s not found", cmd)
}
