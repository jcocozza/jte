package editor

import (
	"bytes"
)

const TAB_STOP = 8

func expandTabs(input []byte) []byte {
	var expanded []byte
	col := 0
	for _, b := range input {
		if b == '\t' {
			spaces := TAB_STOP - (col % TAB_STOP)
			expanded = append(expanded, bytes.Repeat([]byte(" "), spaces)...)
			col += spaces
		} else {
			expanded = append(expanded, b)
			col++
		}
	}
	return expanded
}

type erow struct {
	chars []byte
	render []byte
}

func (r *erow) Render() {
	r.render = expandTabs(r.chars)
}

func (r *erow) InsertChar(at int, c byte) {
	if at < 0 || at > len(r.chars) {
		at = len(r.chars)
	}
	r.chars = append(r.chars, append([]byte{c}, r.chars[at:]...)...)
	r.Render()
}
