package fileutil

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
)

var ErrNoFilename = errors.New("no file name")

// returns the contents, a bool telling you if the file is writeable
func ReadFile(path string) ([][]byte, bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, false, err
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		return nil, false, err
	}
	mode := info.Mode()
	writeable := mode&0200 != 0

	scanner := bufio.NewScanner(f)
	numScans := 0

	contents := [][]byte{}
	for scanner.Scan() {
		line := scanner.Bytes()
		line = bytes.TrimRight(line, "\r\n")
		contents = append(contents, line)
		numScans++
	}
	if numScans == 0 {
		contents = append(contents, []byte{})
	}
	return contents, writeable, nil
}

func Save(filename string, buf []byte) (int, error) {
	if filename == "" {
		return 0, ErrNoFilename
	}
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return 0, fmt.Errorf("unable to open file: %w", err)
	}
	n, err := file.Write(buf)
	if err != nil {
		return 0, fmt.Errorf("unable to write file: %w", err)
	}
	return n, nil
}
