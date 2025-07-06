package panemanager

import (
	"strconv"
	"strings"

	"github.com/jcocozza/jte/internal/buffer"
)

type SplitDirection int

const (
	None       SplitDirection = iota // indicates there is no split recorded
	Horizontal                       // stacked
	Vertical                         // side by side
)

// a doubly linked (complete) tree structure
//
// there should always be 2 children if there are any on a node
type PaneNode struct {
	id        int
	Direction SplitDirection
	Parent    *PaneNode
	First     *PaneNode // First is left or top (depending on direction)
	Second    *PaneNode // Second is right or bottom (depending on direction)
	Ratio     float64   // ratio of first to second
	Bn		  *buffer.BufferNode
}

func newRootPaneNode() *PaneNode {
	return &PaneNode{Ratio: 1}
}

const (
	vsplit = "|"
	hsplit = "-"
)

func (p *PaneNode) isFirst() bool {
	if p.Parent == nil {
		return false
	}
	return p == p.Parent.First
}

func (p *PaneNode) isSecond() bool {
	if p.Parent == nil {
		return false
	}
	return p == p.Parent.Second
}

// solely for debugging
func (p *PaneNode) DrawSimple(spaces int) string {
	idStr := strings.Repeat(" ", spaces) + strconv.Itoa(p.id)
	idStr += "\n"
	if p.IsLeaf() {
		return idStr
	}
	return idStr + p.First.DrawSimple(spaces+1) + p.Second.DrawSimple(spaces+1)
}

// solely for debugging
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
	return ofm(parentId)*10 + 1
}
func secondPaneId(parentId int) int {
	return ofm(parentId)*10 + 2
}

// return the new split node
func (p *PaneNode) splitVertical() *PaneNode {
	p.Direction = Vertical
	p.Ratio = 0.5
	p.First = &PaneNode{id: firstPaneId(p.id), Parent: p}
	p.Second = &PaneNode{id: secondPaneId(p.id), Parent: p}
	return p.First
}

func (p *PaneNode) splitHorizontal() *PaneNode {
	p.Direction = Horizontal
	p.Ratio = 0.5
	p.First = &PaneNode{id: firstPaneId(p.id), Parent: p}
	p.Second = &PaneNode{id: secondPaneId(p.id), Parent: p}
	return p.First
}

func (p *PaneNode) rightMostLeaf() *PaneNode {
	if p.IsLeaf() {
		return p
	}
	return p.Second.rightMostLeaf()
}

func (p *PaneNode) leftMostLeaf() *PaneNode {
	if p.IsLeaf() {
		return p
	}
	return p.First.leftMostLeaf()
}

// return the node to the left of p that is a leaf
func (p *PaneNode) Left() *PaneNode {
	if p.Parent == nil {
		return p
	}
	if p.Parent.Direction == Vertical && p == p.Parent.Second {
		q := p.Parent.First.rightMostLeaf()
		return q
	}
	q := p.Parent.Left()
	if q == p.Parent {
		return p
	}
	return q
}

// return the node to the right of p that is a leaf
func (p *PaneNode) Right() *PaneNode {
	if p.Parent == nil {
		return p
	}
	if p.Parent.Direction == Vertical && p == p.Parent.First {
		q := p.Parent.Second.leftMostLeaf()
		return q
	}
	q := p.Parent.Right()
	if q == p.Parent {
		return p
	}
	return q
}

func (p *PaneNode) bottomMostLeaf() *PaneNode {
	if p.IsLeaf() {
		return p
	}
	return p.Second.bottomMostLeaf()
}

func (p *PaneNode) topMostLeaf() *PaneNode {
	if p.IsLeaf() {
		return p
	}
	return p.First.topMostLeaf()
}

// return the node above p that is a leaf
func (p *PaneNode) Up() *PaneNode {
	if p.Parent == nil {
		return p
	}
	if p.Parent.Direction == Horizontal && p == p.Parent.Second {
		q := p.Parent.First.bottomMostLeaf()
		return q
	}
	q := p.Parent.Up()
	if q == p.Parent {
		return p
	}
	return q
}

// return the node below p that is a leaf
func (p *PaneNode) Down() *PaneNode {
	if p.Parent == nil {
		return p
	}
	if p.Parent.Direction == Horizontal && p == p.Parent.First {
		q := p.Parent.Second.topMostLeaf()
		return q
	}
	q := p.Parent.Down()
	if q == p.Parent {
		return p
	}
	return p
}

// return new current and (optionally) a new root
//
// if the new root is NOT nil, then whereever root is being maintained needs to be updated
func (p *PaneNode) Delete() (*PaneNode, *PaneNode) {
	// we must always have 1 node
	if p.Parent == nil {
		return p, p
	}

	if p.isFirst() {
		sibling := p.Parent.Second
		gp := p.Parent.Parent

		if gp == nil {
			p.Parent = nil
			sibling.Parent = nil
			return sibling, sibling
		}

		switch {
		case p.Parent.isFirst():
			gp.First = sibling
		case p.Parent.isSecond():
			gp.Second = sibling
		default:
			panic("delete")
		}
		sibling.Parent = gp
		return sibling, nil
	}

	if p.isSecond() {
		sibling := p.Parent.First
		gp := p.Parent.Parent
		if gp == nil {
			sibling.Parent = nil
			p.Parent = sibling
			return sibling, sibling
		}
		switch {
		case p.Parent.isFirst():
			gp.First = sibling
		case p.Parent.isSecond():
			gp.Second = sibling
		default:
			panic("delete")
		}
		sibling.Parent = gp
		return sibling, nil
	}
	panic("unexpected state - delete")
}
