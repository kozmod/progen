package exec

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
)

func Test_FileSystemSaveStrategy(t *testing.T) {
	SkipSLowTest(t)

	const (
		pathA = "some_file.txt"
		pathB = "{{ .var }}/some_file.txt"
		pathC = "c/{{ .var }}.txt"

		varTemplateVariableValue = "DATA"
	)

	var (
		dataA = []byte("{{ .var }}")
		dataB = []byte("B")
		dataC = []byte("C")

		templateData = map[string]any{
			"var": varTemplateVariableValue,
		}

		expectedSubPathB = fmt.Sprintf("%s/some_file.txt", varTemplateVariableValue)
		expectedSubPathC = fmt.Sprintf("c/%s.txt", varTemplateVariableValue)
	)
	t.Run("real_save_fs", func(t *testing.T) {
		WithTempDir(t, func(tmpDir string) {
			var (
				tmpDirTarget = filepath.Join(tmpDir, "result")

				expectedPathA = filepath.Join(tmpDirTarget, pathA)
				expectedPathB = filepath.Join(tmpDirTarget, expectedSubPathB)
				expectedPathC = filepath.Join(tmpDirTarget, expectedSubPathC)

				a          = assert.New(t)
				mockLogger = MockLogger{
					infof: func(format string, args ...any) {
						a.NotEmpty(format)
					},
				}
			)

			fs := fstest.MapFS{
				pathA: {
					Data: dataA,
				},
				pathB: {
					Data: dataB,
				},
				pathC: {
					Data: dataC,
				},
			}

			str := NewFileSystemSaveStrategy(fs, templateData, nil, nil, mockLogger)

			dir, err := str.Apply(tmpDirTarget)
			a.NoError(err)
			a.Equal(tmpDirTarget, dir)
			a.FileExists(expectedPathA)
			a.FileExists(expectedPathB)
			a.FileExists(expectedPathC)

			AssertFileDataEqual(t, expectedPathA, []byte(varTemplateVariableValue))
			AssertFileDataEqual(t, expectedPathB, dataB)
			AssertFileDataEqual(t, expectedPathC, dataC)
		})
	})
	t.Run("save_from_real_fs", func(t *testing.T) {
		WithTempDir(t, func(tmpDir string) {
			var (
				pathTempA = filepath.Join(tmpDir, pathA)
				pathTempB = filepath.Join(tmpDir, pathB)
				pathTempC = filepath.Join(tmpDir, pathC)

				a          = assert.New(t)
				mockLogger = MockLogger{
					infof: func(format string, args ...any) {
						a.NotEmpty(format)
					},
				}
				tmpDirTarget = filepath.Join(tmpDir, "result")

				expectedPathA = filepath.Join(tmpDirTarget, pathA)
				expectedPathB = filepath.Join(tmpDirTarget, expectedSubPathB)
				expectedPathC = filepath.Join(tmpDirTarget, expectedSubPathC)
			)

			CreateFile(t, pathTempA, dataA)
			CreateFile(t, pathTempB, dataB)
			CreateFile(t, pathTempC, dataC)

			fs := os.DirFS(tmpDir)

			str := NewFileSystemSaveStrategy(fs, templateData, nil, nil, mockLogger)

			dir, err := str.Apply(tmpDirTarget)
			a.NoError(err)
			a.Equal(tmpDirTarget, dir)
			a.FileExists(expectedPathA)
			a.FileExists(expectedPathB)
			a.FileExists(expectedPathC)

			AssertFileDataEqual(t, expectedPathA, []byte(varTemplateVariableValue))
			AssertFileDataEqual(t, expectedPathB, dataB)
			AssertFileDataEqual(t, expectedPathC, dataC)
		})
	})
}

func Test_DryRunFileSystemSaveStrategy(t *testing.T) {
	const (
		dir = "some_dir"
	)
	mockLogger := MockLogger{
		infof: func(format string, args ...any) {
			assert.NotEmpty(t, format)
			assert.NotEmpty(t, args)
			assert.Equal(t, []any{dir}, args)
		},
	}
	res, err := NewDryRunFileSystemSaveStrategy(mockLogger).Apply(dir)
	assert.NoError(t, err)
	assert.Equal(t, dir, res)
}
