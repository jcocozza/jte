package buffer

// Represents a single row in the buffer
type BufRow []byte

func (b *BufRow) InsertChar(at int, c byte) {
	if at < 0 || at > len(*b) {
		at = len(*b)
	}
	//newChars := make([]byte, len(*b)+1)
	//copy(newChars[:at], (*b)[:at])
	//newChars[at] = c
	//copy(newChars[at+1:], (*b)[at:])
	//*b = newChars
	*b = append((*b)[:at], append([]byte{c}, (*b)[at:]...)...)
}

func (b *BufRow) DelChar(at int) {
	if at < 0 || at >= len(*b) {
		return
	}
	newChars := make([]byte, len(*b)-1)
	copy(newChars[:at], (*b)[:at])
	copy(newChars[at:], (*b)[at+1:])
	*b = newChars
}

func (b *BufRow) append(bytes []byte) {
	*b = append(*b, bytes...)
}

func (b *BufRow) Trim(to int) {
	*b = (*b)[:to]
}
