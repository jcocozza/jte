package editor

import (
	"github.com/jcocozza/jte/internal/buffer"
	"github.com/jcocozza/jte/internal/gutter"
)

type SplitDirection int

const (
	Horizontal SplitDirection = iota // side-by-side
	Vertical                         // top and bottom
)

type Pane struct {
	G   *gutter.Gutter
	Buf *buffer.Buffer
	Active bool // if the cursor is on this node
}

// Pane is nil if this is just a split tracking node
//
// When pane is not nil, we are at a leaf
type SplitNode struct {
	Dir        SplitDirection
	First      *SplitNode
	Second     *SplitNode
	FirstRatio float64 // Ratio of the split

	Pane *Pane // nil for internal
}

func (s *SplitNode) GetUp() {}
func (s *SplitNode) GetDown() {}
func (s *SplitNode) GetLeft() {}
func (s *SplitNode) GetRight() {}

func (s *SplitNode) Resize(newWidth, newHeight int) {
	if s.Dir == Horizontal {
		// Resize horizontally based on ratio
		firstWidth := int(float64(newWidth) * s.FirstRatio)
		secondWidth := newWidth - firstWidth
		s.First.Resize(firstWidth, newHeight)
		s.Second.Resize(secondWidth, newHeight)
	} else {
		// Resize vertically based on ratio
		firstHeight := int(float64(newHeight) * s.FirstRatio)
		secondHeight := newHeight - firstHeight
		s.First.Resize(newWidth, firstHeight)
		s.Second.Resize(newWidth, secondHeight)
	}
}

// return the new "active" split node
func (s *SplitNode) SplitVertical() *SplitNode {
	if s.Pane == nil {
		panic("cannot split an internal node")
	}
	ogPane := s.Pane
	clonedPane := *ogPane
	clonedPane.Active = false
	s.Dir = Vertical
	s.First = &SplitNode{Pane: ogPane}
	s.Second = &SplitNode{Pane: &clonedPane}
	s.FirstRatio = .5
	s.Pane = nil
	return s.First
}

// return the new "active" split node
func (s *SplitNode) SplitHorizontal() *SplitNode {
	if s.Pane == nil {
		panic("cannot split an internal node")
	}
	ogPane := s.Pane
	clonedPane := *ogPane
	clonedPane.Active = false
	s.Dir = Horizontal
	s.First = &SplitNode{Pane: ogPane}
	s.Second = &SplitNode{Pane: &clonedPane}
	s.FirstRatio = .5
	s.Pane = nil
	return s.First
}
