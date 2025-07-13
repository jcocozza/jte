package buffer

type BufRow []rune

func (b *BufRow) Insert(at int, char rune) {
	if at == len(*b) {
		*b = append(*b, char)
		return
	}
	*b = append((*b)[:at], append([]rune{char}, (*b)[at:]...)...)
}

// return the deleted char
func (b *BufRow) Delete(at int) rune {
	if len(*b) == 1 {
		char := (*b)[0]
		*b = []rune{}
		return char
	}
	char := (*b)[at]
	newChars := make([]rune, len(*b)-1)
	copy(newChars[:at], (*b)[:at])
	copy(newChars[at:], (*b)[at+1:])
	*b = newChars
	return char
}
