package renderer

import (
	"log/slog"

	"github.com/jcocozza/jte/internal/editor"
	"github.com/jcocozza/jte/internal/renderer/text"
)

type Renderer interface {
	Setup() error
	Exit(msg string)
	ExitErr(err error)
	Render(e *editor.Editor)
}

func NewRenderer(l *slog.Logger) Renderer {
	return text.NewTextRenderer(l)
}
