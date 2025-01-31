package filesystem

import (
	"fmt"
	"os"
	"regexp"
	"slices"
	"time"
)

// Segment is interface for segment
type Segment interface {
	Write(data []byte) error
	ReadAll() ([][]byte, error)
}

type segment struct {
	file      *os.File
	directory string

	segmentSize    int
	maxSegmentSize int

	fileLib FileLib
}

// NewSegment returns new segment
func NewSegment(directory string, maxSegmentSize int, fileLib FileLib) Segment {
	return &segment{
		directory:      directory,
		maxSegmentSize: maxSegmentSize,
		fileLib:        fileLib,
	}
}

// Write writes bytes of segment
func (s *segment) Write(data []byte) error {
	if s.file == nil || s.segmentSize >= s.maxSegmentSize {
		if err := s.createSegment(); err != nil {
			return fmt.Errorf("failed to create segment file: %w", err)
		}
	}

	writtenBytes, err := s.fileLib.WriteFile(s.file, data)
	if err != nil {
		return fmt.Errorf("failed to write data to segment file: %w", err)
	}

	s.segmentSize += writtenBytes
	return nil
}

func (s *segment) createSegment() error {
	segmentName := fmt.Sprintf("%s/wal_%d.log", s.directory, time.Now().UnixMilli())
	if s.file != nil {
		err := s.file.Close()
		if err != nil {
			return err
		}
	}

	file, err := s.fileLib.CreateFile(segmentName)
	if err != nil {
		return err
	}

	s.file = file
	s.segmentSize = 0
	return nil
}

// ReadAll reads all data from dir
func (s *segment) ReadAll() ([][]byte, error) {
	filenames, err := filenamesFromDir(s.directory)
	if err != nil {
		return nil, err
	}

	return dataFromFiles(s.directory, filenames)
}

func filenamesFromDir(dir string) ([]string, error) {
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

func dataFromFiles(dir string, filenames []string) ([][]byte, error) {
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
