package renderer

import (
	"log/slog"

	"github.com/jcocozza/jte/internal/editor"
)

type LayoutRect struct {
	X, Y       int
	Rows, Cols int
}

type LayoutRenderer struct {
	logger *slog.Logger
}

func NewLayoutRenderer(l *slog.Logger) *LayoutRenderer {
	return &LayoutRenderer{
		logger: l.WithGroup("layout-renderer"),
	}
}

func (r *LayoutRenderer) RenderLayout(e *editor.Editor, root *editor.SplitNode, pr PaneRenderer, screenrows int, screencols int) [][]byte {
	screen := make([][]byte, screenrows)
	for i := range screen {
		screen[i] = make([]byte, screencols)
		for j := range screen[i] {
			screen[i][j] = ' '
		}
	}
	r.RenderNode(e, root, pr, LayoutRect{X: 0, Y: 0, Rows: screenrows, Cols: screencols}, screen)
	return screen
}

func (r *LayoutRenderer) RenderNode(e *editor.Editor, node *editor.SplitNode, pr PaneRenderer, rect LayoutRect, screen [][]byte) {
	if node == nil {
		return
	}
	if node.Pane != nil {
		psd := PaneStatusData{ Active: node.Pane.Active, Mode: e.Mode() }
		rendered := pr.Render(rect.Rows, rect.Cols, psd, node.Pane.G, node.Pane.Buf)
		for i := 0; i < len(rendered) && i+rect.Y < len(screen); i++ {
			copy(screen[i+rect.Y][rect.X:], rendered[i])
		}
		return
	}
	if node.Dir == editor.Vertical {
		firstW := int(float64(rect.Cols) * node.FirstRatio)
		secondW := rect.Cols - firstW
		r.RenderNode(e, node.First, pr, LayoutRect{rect.X, rect.Y, rect.Rows, firstW}, screen)
		r.RenderNode(e, node.Second, pr, LayoutRect{rect.X + firstW, rect.Y, rect.Rows, secondW}, screen)
	} else {
		firstH := int(float64(rect.Rows) * node.FirstRatio)
		secondH := rect.Rows - firstH
		r.RenderNode(e, node.First, pr, LayoutRect{rect.X, rect.Y, firstH, rect.Cols}, screen)
		r.RenderNode(e, node.Second, pr, LayoutRect{rect.X, rect.Y + firstH, secondH, rect.Cols}, screen)
	}
}
