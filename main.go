package main

import (
	"os"
	"time"

	"github.com/jcocozza/jte/editor"
)

func main() {
	e := editor.InitEditor()
	defer e.Exit("regular exit")
	if len(os.Args) >= 2 {
		err := e.Open(os.Args[1])
		if err != nil {
			e.ExitErr(err)
		}
	}
	e.SetMsg("Momento Mori", time.Duration(3*time.Second))
	//e.Debug()
	for {
		e.Refresh()
		e.ProcessKeypress()
	}
}
