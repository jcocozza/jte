package renderer

import "os"

type abuf []byte

func (b *abuf) Append(bytes []byte) {
	*b = append(*b, bytes...)
}

func (b *abuf) Clear() {
	*b = []byte{}
}

func (b *abuf) Flush() {
	os.Stdout.Write(*b)
	b.Clear()
}
