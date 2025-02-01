package wal

import (
	"os"
	"testing"

	"concurrency_go_course/internal/filesystem"
	"concurrency_go_course/pkg/logger"

	"github.com/stretchr/testify/assert"
)

const (
	testDataDir     = "tmp"
	testDataDirRead = "read_tmp"
)

func TestLogsManagerWrite(t *testing.T) {
	logger.MockLogger()

	if _, err := os.Stat(testDataDir); err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(testDataDir, os.ModePerm) //nolint:gosec
			if err != nil {
				t.Errorf("unable to create temp dir: mkdir error: %s", err)
			}
		} else {
			t.Errorf("unable to get dir info: %s", err)
		}
	}

	defer func() {
		err := os.RemoveAll(testDataDir)
		if err != nil {
			t.Errorf("unable to remove tmp dir for test data [%s]: %s", testDataDir, err)
		}
	}()

	req1 := Request{
		Command:    "SET",
		Args:       []string{"key", "value"},
		doneStatus: make(chan error, 1),
	}

	requests := []Request{
		req1,
	}

	fileLib := filesystem.NewFileLib()
	segment := filesystem.NewSegment(testDataDir, 10, fileLib)

	logsManager, err := NewLogsManager(segment)
	if err != nil {
		t.Errorf("failed: %s", err)
	}

	logsManager.Write(requests)

	for _, request := range requests {
		err := <-request.doneStatus
		if err != nil {
			t.Errorf("failed: %s", err)
		}

		_, ok := <-request.doneStatus
		if ok {
			t.Errorf("failed: channel was not closed [request: %+v]", request)
		}
	}
}

func TestLogsManagerReadAll(t *testing.T) {
	t.Parallel()
	logger.MockLogger()

	err := os.Mkdir(testDataDirRead, os.ModePerm) //nolint:gosec
	if err != nil {
		t.Errorf("cannot create temporary dir for test data [%s]: %s", testDataDirRead, err)
	}

	defer func() {
		err := os.RemoveAll(testDataDirRead)
		if err != nil {
			t.Errorf("unable to remove tmp dir for test data [%s]: %s", testDataDirRead, err)
		}
	}()

	requests := []Request{
		{
			Command:    "SET",
			Args:       []string{"key", "value"},
			doneStatus: make(chan error, 1),
		},
		{
			Command:    "DEL",
			Args:       []string{"key1"},
			doneStatus: make(chan error, 1),
		},
	}

	fileLib := filesystem.NewMockFileLib()
	segment := filesystem.NewSegment(testDataDirRead, 10, fileLib)

	logsManager, err := NewLogsManager(segment)
	if err != nil {
		t.Errorf("failed: %s", err)
	}

	logsManager.Write(requests)

	segmentR := filesystem.NewSegment(testDataDirRead, 10, fileLib)

	logsManager, err = NewLogsManager(segmentR)
	if err != nil {
		t.Errorf("failed: %s", err)
	}

	reqs, err := logsManager.ReadAll()
	if err != nil {
		t.Errorf("unable to read segments data from [%s]", testDataDirRead)
	}

	for _, r := range reqs {
		assert.Equal(t, r.Command, r.Command)
		assert.Equal(t, r.Args, r.Args)
	}
}
