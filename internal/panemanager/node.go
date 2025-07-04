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

// a doubly linked (complete) tree structure
//
// there should always be 2 children if there are any on a node
type PaneNode struct {
	id        int
	Direction SplitDirection
	Parent    *PaneNode
	// First is left or top
	First *PaneNode
	// Second is right or bottom
	Second *PaneNode
	Ratio  float64 // ratio of first to second
	Active bool
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
	return ofm(parentId)*10 + 1
}
func secondPaneId(parentId int) int {
	return ofm(parentId)*10 + 2
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
// and set it to active
func (p *PaneNode) Left() *PaneNode {
	if p.Parent == nil {
		return p
	}
	if p.Parent.Direction == Vertical && p == p.Parent.Second {
		q := p.Parent.First.rightMostLeaf()
		p.Active = false
		q.Active = true
		return q
	}
	q := p.Parent.Left()
	if q == p.Parent {
		return p
	}
	p.Active = false
	q.Active = true
	return q
}

// return the node to the right of p that is a leaf
// and set it to active
func (p *PaneNode) Right() *PaneNode {
	if p.Parent == nil {
		return p
	}
	if p.Parent.Direction == Vertical && p == p.Parent.First {
		q := p.Parent.Second.leftMostLeaf()
		p.Active = false
		q.Active = true
		return q
	}
	q := p.Parent.Right()
	if q == p.Parent {
		return p
	}
	p.Active = false
	q.Active = true
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
// and set it to active
func (p *PaneNode) Up() *PaneNode {
	if p.Parent == nil {
		return p
	}
	if p.Parent.Direction == Horizontal && p == p.Parent.Second {
		q := p.Parent.First.bottomMostLeaf()
		p.Active = false
		q.Active = true
		return q
	}
	q := p.Parent.Up()
	if q == p.Parent {
		return p
	}
	p.Active = false
	q.Active = true
	return q
}

// return the node below p that is a leaf
// and set it to active
func (p *PaneNode) Down() *PaneNode {
	if p.Parent == nil {
		return p
	}
	if p.Parent.Direction == Horizontal && p == p.Parent.First {
		q := p.Parent.Second.topMostLeaf()
		p.Active = false
		q.Active = true
		return q
	}
	q := p.Parent.Down()
	if q == p.Parent {
		return p
	}
	p.Active = false
	q.Active = true
	return p
}

// move n1 to n2's position and drop n2
//
// assume n2 is parent of n1 and has a parent of its own
func swapAndDrop(n1 *PaneNode, n2 *PaneNode) *PaneNode {
	gp := n2.Parent
	n1.Parent = gp
	n1.Active = true
	switch {
	case n2.isFirst():
		gp.First = n1
	case n2.isSecond():
		gp.Second = n1
	default:
		panic("unexpected state - swap and drop")
	}
	return n1
}

func (p *PaneNode) Delete() *PaneNode {
	p.Active = false

	// we should always have 1 node
	if p.Parent == nil {
		return p
	}
	// the parent is root, so alternate becomes root
	// TODO: this case doesn't seem to be working right
	if p.Parent.Parent == nil {
		var promotee *PaneNode
		if p.isFirst() {
			promotee = p.Parent.Second
			p.Parent.First = nil
		} else if p.isSecond() {
			promotee = p.Parent.First
			p.Parent.Second = nil
		} else {
			panic("unexpected state")
		}
		promotee.Parent = nil
		promotee.Active = true
		return promotee
	}

	// now we are guaranteed parent and grandparent nodes
	switch {
	case p.isFirst():
		return swapAndDrop(p.Parent.Second, p.Parent)
	case p.isSecond():
		return swapAndDrop(p.Parent.First, p.Parent)
	default:
		panic("unexpected state")
	}
}
