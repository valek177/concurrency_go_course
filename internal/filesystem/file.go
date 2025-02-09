package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
)

// FileLib is interface for file management lib
type FileLib interface {
	CreateFile(filename string) (*os.File, error)
	WriteFile(file *os.File, data []byte) (int, error)
	DataFromFiles(dir string, filenames []string) ([][]byte, error)
	FilenamesFromDir(dir string) ([]string, error)
	SegmentNext(dir, filename string) (string, error)
	SegmentLast(dir string) (string, error)
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

// DataFromFiles returns data from files
func (f *filelib) DataFromFiles(dir string, filenames []string) ([][]byte, error) {
	dataRes := make([][]byte, 0, len(filenames))

	for _, f := range filenames {
		data, err := os.ReadFile(fmt.Sprintf("%s/%s", dir, f))
		if err != nil {
			return nil, err
		}

		dataRes = append(dataRes, data)
	}

	return dataRes, nil
}

// FilenamesFromDir returns list of files from dir
func (f *filelib) FilenamesFromDir(dir string) ([]string, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("unable to read WAL directory: %w", err)
	}

	fileNames := make([]string, 0, len(files))
	re := regexp.MustCompile(`wal_\d+\.log`)

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if !re.MatchString(file.Name()) {
			continue
		}
		fileNames = append(fileNames, file.Name())
	}

	slices.Sort(fileNames)

	return fileNames, nil
}

func (f *filelib) SegmentLast(dir string) (string, error) {
	wals, err := f.FilenamesFromDir(dir)
	if err != nil {
		return "", err
	}
	if len(wals) == 0 {
		return "", fmt.Errorf("no segments found")
	}
	return wals[len(wals)-1], nil
}

func (f *filelib) SegmentNext(dir, filename string) (string, error) {
	wals, err := f.FilenamesFromDir(dir)
	if err != nil {
		return "", err
	}

	// Get newest WAL
	for i := len(wals) - 1; i >= 0; i-- {
		if wals[i] > filename {
			return wals[i], nil
		}
	}

	return "", fmt.Errorf("unable to find next segment")
}
