package exec

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kozmod/progen/internal/entity"
)

const (
	tempDir = ""
)

// WithTempDir create a temporary directory for testing, test function and remove a temporary directory after test execution
func WithTempDir(t *testing.T, test func(dir string)) {
	t.Helper()

	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	tmpPath, err := os.MkdirTemp(tempDir, fmt.Sprintf("%s_", filepath.Base(f.Name())))
	assert.NoError(t, err)
	defer func() {
		err = os.RemoveAll(tmpPath)
		if err != nil {
			t.Fatalf("remove all in tmp dir: %v", err)
		}
	}()
	test(tmpPath)
}

func CreateFile(t *testing.T, path string, data []byte) {
	t.Helper()

	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, os.ModePerm)
		assert.NoError(t, err)
	}

	//nolint:gosec
	err := os.WriteFile(path, data, os.ModePerm)
	assert.NoError(t, err)
}

func SkipSLowTest(t *testing.T) {
	t.Helper()

	if testing.Short() {
		pc := make([]uintptr, 1)
		runtime.Callers(2, pc)
		f := runtime.FuncForPC(pc[0])
		t.Skipf("skipping slow test: %s", f.Name())
	}
}

func AssertFileDataEqual(t *testing.T, path string, expData []byte) {
	t.Helper()

	data, err := os.ReadFile(path)
	assert.NoError(t, err)
	assert.Equal(t, expData, data)
}

type MockLogger struct {
	entity.Logger
	infof func(format string, args ...any)
}

func (m MockLogger) Infof(format string, args ...any) {
	m.infof(format, args...)
}

type MockProducer struct {
	file entity.DataFile
	err  error
}

func (m *MockProducer) Get() (entity.DataFile, error) {
	if m.err != nil {
		return entity.DataFile{}, m.err
	}
	return m.file, nil
}

type MockExecutor struct{}

func (m MockExecutor) Exec() error {
	return nil
}

type MockFileProc struct{}

func (m MockFileProc) Apply(_ entity.DataFile) (entity.DataFile, error) {
	return entity.DataFile{}, nil
}
