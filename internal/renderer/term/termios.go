package term

type Term interface {
	WindowSize() (int, int, error)
}

type RawMode struct {
	originalState any
	fd            int
}

func EnableRawMode() (*RawMode, error) {
	return enableRawMode()
}

func (r *RawMode) Restore() error {
	return restore(r)
}

func (r *RawMode) WindowSize() (int, int, error) {
	return getWindowSize()
}
