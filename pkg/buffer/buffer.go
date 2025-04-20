package buffer

type FileType string

// cursor location in the buffer
// 	X - row
//  Y - column
type Cursor struct {
	X int
	Y int
}

// represents a single row
type row []byte

type Buffer struct {
	// a unique identifier
	uuid string

	// the rows in the underlying file
	rows []row
	cursor Cursor

	// state stuff
	Modified bool
	ReadOnly bool
	FileType FileType

	// file stuff
	FileName string
	FilePath string
}
