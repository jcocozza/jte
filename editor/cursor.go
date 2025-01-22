package editor

// the cursor in the terminal
//
// X and Y cannot be less then 0
type cursor struct {
	X int
	Y int
	// rendered X
	// RX int
}
