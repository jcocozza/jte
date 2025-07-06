package panemanager

import (
	"fmt"
	"testing"
)

func PM() *PaneManager {
	r := ComplexSplit()

	return &PaneManager{
		Root: r,
		Curr: r.First,
	}
}

func TestMovement(t *testing.T) {
	var tests = []struct{
		name string
		pm *PaneManager
	}{
		{name: "complex", pm: PM()},
	}

	for _,tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.pm.Right()
			tt.pm.Right()
			//tt.p.First.Right().Right().Delete().Delete().Delete()
			//p.Delete()//.Delete()//.Delete()
			fmt.Println(tt.pm.Root.DrawSimple(0))
			fmt.Println(tt.pm.Curr.id)
		})
	}
}

func TestMovementAndDelete(t *testing.T) {
	var tests = []struct{
		name string
		pm *PaneManager
	}{
		{name: "complex", pm: PM()},
	}

	for _,tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.pm.Right()
			tt.pm.Right()
			tt.pm.Delete()
			tt.pm.Delete()
			tt.pm.Delete()
			fmt.Println(tt.pm.Root.DrawSimple(0))
			fmt.Println(tt.pm.Root)
			fmt.Println(tt.pm.Curr.id)
		})
	}
}
