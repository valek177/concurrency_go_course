package filesystem

import (
	"os"
	"path/filepath"
)

// FileLib is interface for file management lib
type FileLib interface {
	CreateFile(filename string) (*os.File, error)
	WriteFile(file *os.File, data []byte) (int, error)
}

type filelib struct{}

// NewFileLib returns new FileLib
func NewFileLib() FileLib {
	var filelib filelib
	return &filelib
}

// CreateFile creates new file
func (f *filelib) CreateFile(filename string) (*os.File, error) {
	file, err := os.OpenFile(filepath.Clean(filename), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, //nolint:gosec
		os.ModePerm)
	if err != nil {
		return nil, err
	}

	return file, err
}

// WriteFile writes data to file by file descriptor
func (f *filelib) WriteFile(file *os.File, data []byte) (int, error) {
	writtenBytes, err := file.Write(data)
	if err != nil {
		return 0, err
	}

	if err = file.Sync(); err != nil {
		return 0, err
	}

	return writtenBytes, nil
}
