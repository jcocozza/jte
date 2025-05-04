package renderer

type Renderer interface {
	Setup() error
	Exit(msg string)
	ExitErr(err error)
	Render()
}
