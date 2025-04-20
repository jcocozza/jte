package actions


type Action interface {
	Name() string
	Apply()
}

// implements the Action interface
type ActionFunc struct {
	name  string
	apply func()
}

func NewAction(name string, fn func()) Action {
	return &ActionFunc{name: name, apply: fn}
}

func (a *ActionFunc) Name() string { return a.name }

func (a *ActionFunc) Apply() { a.apply() }

// define actions that can be taken

var (
	CursorUp    = NewAction("cursor up", func() {})
	CursorDown  = NewAction("cursor down", func() {})
	CursorLeft  = NewAction("cursor left", func() {})
	CursorRight = NewAction("cursor right", func() {})
	Exit        = NewAction("exit", func() {})
)
