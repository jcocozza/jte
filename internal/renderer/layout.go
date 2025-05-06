package renderer

import (
	"bytes"

	"github.com/jcocozza/jte/internal/editor"
)

type LayoutRect struct {
	X, Y       int
	Rows, Cols int
}

func RenderLayout(root *editor.SplitNode, pr PaneRenderer, screenrows int, screencols int) [][]byte {
	screen := make([][]byte, screenrows)
	for i := range screen {
		screen[i] = make([]byte, screencols)
		for j := range screen[i] {
			screen[i][j] = ' '
		}
	}
	RenderNode(root, pr, LayoutRect{X: 0, Y:0, Rows: screenrows, Cols: screencols}, screen)
	return screen
}

func RenderNode(node *editor.SplitNode, pr PaneRenderer, rect LayoutRect, screen [][]byte) {
	if node == nil {
		return
	}
	if node.Pane != nil {
		rendered := pr.Render(rect.Rows, rect.Cols, node.Pane.G, node.Pane.Buf)
		for i := 0; i < len(rendered) && i+rect.Y < len(screen); i++ {
			copy(screen[i+rect.Y][rect.X:], rendered[i])
		}
		return
	}

	if node.Dir == editor.Horizontal {
		firstW := int(float64(rect.Cols) * node.FirstRatio)
		secondW := rect.Cols - firstW
		RenderNode(node.First, pr, LayoutRect{rect.X, rect.Y, firstW, rect.Rows}, screen)
		RenderNode(node.Second, pr, LayoutRect{rect.X + firstW, rect.Y, secondW, rect.Rows}, screen)
	} else {
		firstH := int(float64(rect.Rows) * node.FirstRatio)
		secondH := rect.Rows - firstH
		RenderNode(node.First, pr, LayoutRect{rect.X, rect.Y, rect.Cols, firstH}, screen)
		RenderNode(node.Second, pr, LayoutRect{rect.X, rect.Y + firstH, rect.Cols, secondH}, screen)
	}
}

func Render(width, height int, sn *editor.SplitNode) [][]byte {
	var firstLines, secondLines [][]byte

	if sn.Dir == editor.Horizontal {
		w1 := int(float64(width) * sn.FirstRatio)
		w2 := width - w1
		firstLines = Render(w1, height, sn.First)
		secondLines = Render(w2, height, sn.Second)
		return mergeHorizontal(firstLines, secondLines, w1, w2)
	} else {
		h1 := int(float64(height) * sn.FirstRatio)
		h2 := height - h1
		firstLines = Render(width, h1, sn.First)
		secondLines = Render(width, h2, sn.Second)
		return append(firstLines, secondLines...)
	}
}

func mergeHorizontal(left, right [][]byte, w1, w2 int) [][]byte {
	maxLines := max(len(left), len(right))
	result := make([][]byte, maxLines)

	for i := 0; i < maxLines; i++ {
		l := getLine(left, i)
		r := getLine(right, i)
		result[i] = append(padToWidth(l, w1), padToWidth(r, w2)...)
	}
	return result
}

func getLine(lines [][]byte, i int) []byte {
	if i < len(lines) {
		return lines[i]
	}
	return []byte{}
}

func padToWidth(s []byte, width int) []byte {
	if len(s) > width {
		return s[:width]
	}
	return append(s, bytes.Repeat([]byte(" "), width-len(s))...)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
