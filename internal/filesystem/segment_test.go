package filesystem

import (
	"os"
	"strings"
	"testing"

	"concurrency_go_course/pkg/logger"
)

const (
	testDataDir     = "tmp"
	testDataDirRead = "test_data"
)

func TestSegmentWrite(t *testing.T) {
	t.Parallel()

	err := os.Mkdir(testDataDir, os.ModePerm) //nolint:gosec
	if err != nil {
		t.Errorf("cannot create temporary dir for test data [%s]: %s", testDataDir, err)
	}

	defer func() {
		err := os.RemoveAll(testDataDir)
		if err != nil {
			t.Errorf("unable to remove tmp dir for test data [%s]: %s", testDataDir, err)
		}
	}()

	mockFileLib := NewMockFileLib()

	segment := NewSegment(testDataDir, 10, mockFileLib)

	err = segment.Write([]byte("aaaaa"))
	if err != nil {
		t.Errorf("unable to write test data: %s", err)
	}

	err = segment.Write([]byte("bbbbb"))
	if err != nil {
		t.Errorf("unable to write test data: %s", err)
	}

	stat, err := os.Stat(testDataDir + "/wal_1.log")
	if err != nil {
		t.Errorf("unable to get file info [%s]", testDataDir+"/wal_1.log")
	}

	if stat.Size() != 10 {
		t.Errorf("wrong file size: expected 10, got %d", stat.Size())
	}
}

func TestSegmentReadAll(t *testing.T) {
	t.Parallel()
	logger.MockLogger()

	mockFileLib := NewMockFileLib()
	segment := NewSegment(testDataDirRead, 10, mockFileLib)

	data, err := segment.ReadAll()
	if err != nil {
		t.Errorf("unable to read segments data from [%s]", testDataDirRead)
	}

	if len(data) != 2 {
		t.Errorf("wrong number of segments: expected %d, got %d", 2, len(data))
	}

	if strings.TrimSuffix(string(data[0]), "\n") != "SET k v" {
		t.Errorf("wrong segment data: expected %s, got %s", "'SET k v'", string(data[0]))
	}

	if strings.TrimSuffix(string(data[1]), "\n") != "SET k1 v1" {
		t.Errorf("wrong segment data: expected %s, got %s", "'SET k1 v1'", string(data[1]))
	}
}
