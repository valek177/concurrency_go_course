package filesystem

import "os"

// MockFileLib is interface for file lib
type MockFileLib interface {
	CreateFile(filename string) (*os.File, error)
	WriteFile(file *os.File, data []byte) (int, error)
}

type mockfilelib struct{}

// NewMockFileLib mocks file lib
func NewMockFileLib() MockFileLib {
	var lib mockfilelib
	return &lib
}

// CreateFile creates new file
func (f *mockfilelib) CreateFile(_ string) (*os.File, error) {
	file, err := os.OpenFile("tmp/wal_1.log", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, //nolint:gosec
		os.ModePerm)
	if err != nil {
		return nil, err
	}

	return file, err
}

// WriteFile writes data to file by file descriptor
func (f *mockfilelib) WriteFile(file *os.File, data []byte) (int, error) {
	writtenBytes, err := file.Write(data)
	if err != nil {
		return 0, err
	}

	if err = file.Sync(); err != nil {
		return 0, err
	}

	return writtenBytes, nil
}
