package panemanager

import (
	"strconv"
	"strings"
)

type SplitDirection int

const (
	None       SplitDirection = iota // indicates there is no split recorded
	Horizontal                       // stacked
	Vertical                         // side by side
)

// a doubly linked tree structure
type PaneNode struct {
	// each order of magnitide increase in the id represents a tree level
	id        int
	Direction SplitDirection
	Parent    *PaneNode
	First     *PaneNode
	Second    *PaneNode
	Ratio     float64 // ratio of first to second
	Active    bool
}

const (
	vsplit = "|"
	hsplit = "-"
)

// solely for debugging
func (p *PaneNode) DrawSimple(spaces int) string {
	idStr := strings.Repeat(" ", spaces) + strconv.Itoa(p.id)
	if p.Active {
		idStr += "(A)"
	}
	idStr += "\n"
	if p.IsLeaf() {
		return idStr
	}
	return idStr + p.First.DrawSimple(spaces+1) + p.Second.DrawSimple(spaces+1)
}

func (p *PaneNode) Draw(s [][]string, startX int, startY int, rows int, cols int) [][]string {
	switch p.Direction {
	case Vertical:
		firstW := int(float64(cols) * p.Ratio)
		secondW := cols - firstW
		for i := range rows {
			s[i+startY][firstW+startX] = vsplit
		}
		p.First.Draw(s, startX, startY, rows, firstW)
		p.Second.Draw(s, startX+firstW, startY, rows, secondW)
	case Horizontal:
		firstH := int(float64(rows) * p.Ratio)
		secondH := rows - firstH
		for i := range cols {
			s[firstH+startY][i+startX] = hsplit
		}
		p.First.Draw(s, startX, startY, rows-firstH, cols)
		p.Second.Draw(s, startX, startY+firstH, secondH, cols)
	case None:
		s[startY+3][startX+3] = strconv.Itoa(p.id)
	}
	return s
}

// solely for debugging
func (p *PaneNode) String() string {
	width, height := 100, 25
	s := make([][]string, height)
	for i := range s {
		s[i] = make([]string, width)
	}
	s = p.Draw(s, 0, 0, height, width)

	str := ""

	for _, row := range s {
		for _, elm := range row {
			if elm != "" {
				str += elm
			} else {
				str += " "
			}
		}
		str += "\n"
	}
	return str
}

func (p *PaneNode) IsLeaf() bool {
	return p.First == nil && p.Second == nil
}

// this currently does nothing
func (p *PaneNode) Resize(newWidth, newHeight int) {
	if p.Direction == Horizontal {
		// Resize horizontally based on ratio
		firstWidth := int(float64(newWidth) * p.Ratio)
		secondWidth := newWidth - firstWidth
		p.First.Resize(firstWidth, newHeight)
		p.Second.Resize(secondWidth, newHeight)
	} else {
		// Resize vertically based on ratio
		firstHeight := int(float64(newHeight) * p.Ratio)
		secondHeight := newHeight - firstHeight
		p.First.Resize(newWidth, firstHeight)
		p.Second.Resize(newWidth, secondHeight)
	}
}

func ofm(n int) int {
	if n == 0 {
		return 0
	}
	res := 1
	for n >= 10 {
		n /= 10
		res *= 10
	}
	return res
}

func firstPaneId(parentId int) int {
	return ofm(parentId) * 10 + 1
}
func secondPaneId(parentId int) int {
	return ofm(parentId) * 10 + 2
}

// return the new split node
func (p *PaneNode) SplitVertical() *PaneNode {
	p.Direction = Vertical
	p.Ratio = 0.5
	p.First = &PaneNode{id: firstPaneId(p.id), Parent: p}
	p.Second = &PaneNode{id: secondPaneId(p.id), Parent: p}
	return p.First
}

func (p *PaneNode) SplitHorizontal() *PaneNode {
	p.Direction = Horizontal
	p.Ratio = 0.5
	p.First = &PaneNode{id: firstPaneId(p.id), Parent: p}
	p.Second = &PaneNode{id: secondPaneId(p.id), Parent: p}
	return p.First
}

// TODO
// return the node to the left of p that is a leaf
func (p *PaneNode) Left() *PaneNode { return p }

// TODO
// return the node to the right of p that is a leaf
func (p *PaneNode) Right() *PaneNode { return p }

// TODO
// return the node above p that is a leaf
func (p *PaneNode) Up() *PaneNode { return p }

// TODO
// return the node below p that is a leaf
func (p *PaneNode) Down() *PaneNode { return p }
