package main

import (
	"os"
	"unnamed/editor"
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
	e.SetMsg("Momento Mori", -1)
	//e.Debug()
	for {
		e.Refresh()
		e.ProcessKeypress()
	}
}
