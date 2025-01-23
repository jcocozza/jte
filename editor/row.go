package editor

import (
	"bytes"

	"github.com/jcocozza/jte/color"
)

const TAB_STOP = 8

func expandTabs(input []byte) []byte {
	var expanded []byte
	col := 0
	var currHL color.Highlight = -1
	for _, b := range input {
		if b == '\t' {
			spaces := TAB_STOP - (col % TAB_STOP)
			expanded = append(expanded, bytes.Repeat([]byte(" "), spaces)...)
			col += spaces
		} else {
			cbh := color.ColorByte(b)
			switch cbh {
			case color.HL_NORMAL:
				if currHL != -1 {
					expanded = append(expanded, color.MakeColor(color.RESET)...)
					currHL = -1
				}
				expanded = append(expanded, b)
			default:
				colr := color.SyntaxToColor(cbh)
				if cbh != currHL {
					currHL = cbh
					expanded = append(expanded, colr...)
				}
				expanded = append(expanded, b)
			}
			col++
		}
	}
	if currHL != -1 {
		expanded = append(expanded, color.MakeColor(color.RESET)...)
	}
	return expanded
}

type erow struct {
	chars  []byte
	render []byte
	//hl     []rune
}

func (r *erow) Render() {
	r.render = expandTabs(r.chars)
	//r.highlight()
}

/*
func (r *erow) highlight() {
	r.hl = make([]rune, len(r.render))
	for i, c := range r.render {
		if unicode.IsDigit(rune(c)) {
			r.hl[i] = color.HL_NUMBER
		}
	}
}
*/

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
