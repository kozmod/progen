package exec

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_MkdirAllProc(t *testing.T) {
	SkipSLowTest(t)

	const (
		someDir = "some_dir"
	)

	WithTempDir(t, func(tmpDir string) {
		var (
			exp        = filepath.Join(tmpDir, someDir)
			mockLogger = MockLogger{
				infof: func(format string, args ...any) {
					assert.NotEmpty(t, format)
					assert.ElementsMatch(t, []string{exp}, args)
				},
			}
		)

		res, err := NewMkdirAllProc(mockLogger).Process(exp)
		assert.NoError(t, err)
		assert.Equal(t, exp, res)
		assert.DirExists(t, res)
	})
}
