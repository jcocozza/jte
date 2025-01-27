package fileutil

import (
	"errors"
	"fmt"
	"os"
)

var ErrNoFilename = errors.New("no file name")

func OpenOrCreateFile(filename string) (*os.File, error) {
	if filename == "" {
		return os.CreateTemp("", "jte_*")
	}
	return os.Open(filename)
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
