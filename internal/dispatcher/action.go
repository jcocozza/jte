package dispatcher

type Action interface {
	Apply()
}
