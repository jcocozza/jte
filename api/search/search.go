package search

import (
	"bytes"

	"github.com/jcocozza/jte/api/buffer"
)

type Location struct {
	X int
	Y int
}

func findAllIndicies(pattern string, row []byte) []int {
	var indices []int
	start := 0
	for {
		index := bytes.Index(row[start:], []byte(pattern))
		if index == -1 {
			break
		}
		indices = append(indices, start+index)
		start += index + 1
	}
	return indices
}

// a very basic search
//
// simply checks all rows in the buffer for the pattern
func Search(pattern string, buf buffer.Buffer) []Location {
	locs := []Location{}
	for y := range buf.NumRows() {
		row := buf.Row(y)
		indices := findAllIndicies(pattern, row)
		for _, idx := range indices {
			locs = append(locs, Location{X: idx, Y:y})
		}
	}
	return locs
}
