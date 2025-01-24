package editor

import (
	"bytes"

	"github.com/jcocozza/jte/color"
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
	chars  []byte
	render []byte
	hl     []color.Highlight
}

func (r *erow) Render() {
	r.render = expandTabs(r.chars)
	r.setHighlight()
}

func (r *erow) setHighlight() {
	r.hl = make([]color.Highlight, len(r.render))
	for i, c := range r.render {
		r.hl[i] = color.ColorByte(c)
	}
}

func (r *erow) InsertChar(at int, c byte) {
	if at < 0 || at > len(r.chars) {
		at = len(r.chars)
	}
	newChars := make([]byte, len(r.chars)+1)
	copy(newChars[:at], r.chars[:at])
	newChars[at] = c
	copy(newChars[at+1:], r.chars[at:])
	r.chars = newChars
	r.Render()
}

func (r *erow) DelChar(at int) {
	if at < 0 || at >= len(r.chars) {
		return
	}
	newChars := make([]byte, len(r.chars)-1)
	copy(newChars[:at], r.chars[:at])
	copy(newChars[at:], r.chars[at+1:])
	r.chars = newChars
	r.Render()
}

func (r *erow) append(bytes []byte) {
	r.chars = append(r.chars, bytes...)
	r.Render()
}

func (r *erow) Trim(to int) {
	r.chars = r.chars[:to]
	r.Render()
}
