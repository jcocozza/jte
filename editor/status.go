package editor

import (
	"fmt"
)

type status struct {
	filename string
}

func (s *status) Bar() string {
	var displayName string = s.filename
	if s.filename == "" {
		displayName = "[No Name]"
	}
	return fmt.Sprintf("%f - %d lines %s", displayName)
}
