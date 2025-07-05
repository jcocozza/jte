package panemanager

// go test -v .

import (
	"fmt"
	"testing"
)

var root = &PaneNode{
	Ratio: 1,
}

func VSplit() *PaneNode {
	r := &PaneNode{ Ratio: 1}
	r.SplitVertical()
	return r
}

func HSplit() *PaneNode {
	r := &PaneNode{ Ratio: 1}
	r.SplitHorizontal()
	return r
}

func ComplexSplit() *PaneNode {
	r := &PaneNode{ Ratio: 1}
	r.SplitVertical()
	r.Second.SplitHorizontal().SplitVertical()
	return r
}

func TestDraw(t *testing.T) {
	var tests = []struct {
		name string
		p    *PaneNode
	}{
		{name: "root", p: root},
		{name: "vsplit", p: VSplit()},
		{name: "hsplit", p: HSplit()},
		{name: "complex split", p: ComplexSplit()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Println(tt.p.DrawSimple(0))
			fmt.Println(tt.p)
		})
	}
}
