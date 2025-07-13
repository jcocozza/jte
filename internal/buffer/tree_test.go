package buffer

import (
	"fmt"
	"testing"
)

var root = &ChangeNode{
	Children: []*ChangeNode{
		{
			Children: []*ChangeNode{
				{
					Children:  []*ChangeNode{{}},
				},
			},
		},
		{
			Children: []*ChangeNode{{},{ Children: []*ChangeNode{{}}},{}},
		},
	},
}

func TestDraw(t *testing.T) {
	var tests = []struct {
		name string
		n *ChangeNode
	}{
		{name: "root", n: root},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Println(tt.n)
		})
	}
}
