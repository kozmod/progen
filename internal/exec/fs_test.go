package exec

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kozmod/progen/internal/entity"
)

func Test_FileSystemStrategy(t *testing.T) {
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

	t.Run("stub", func(t *testing.T) {
		WithTempDir(t, func(tmpDir string) {
			var (
				pathTempA = filepath.Join(tmpDir, pathA)
				pathTempB = filepath.Join(tmpDir, pathB)
				pathTempC = filepath.Join(tmpDir, pathC)

				expectedPathB = filepath.Join(tmpDir, expectedSubPathB)
				expectedPathC = filepath.Join(tmpDir, expectedSubPathC)

				a          = assert.New(t)
				mockLogger = MockLogger{
					infof: func(format string, args ...any) {
						a.NotEmpty(format)
					},
				}
			)

			CreateFile(t, pathTempA, dataA)
			CreateFile(t, pathTempB, dataB)
			CreateFile(t, pathTempC, dataC)

			str := FileSystemStrategy{
				templateProcFn: func() entity.TemplateProc {
					return entity.NewTemplateProc(templateData, nil, nil)
				},
				dirExecutorFn: func(dirs []string) entity.Executor {
					a.ElementsMatch([]string{filepath.Dir(expectedPathB), filepath.Dir(expectedPathC)}, dirs)
					return MockExecutor{}
				},
				strategiesFn: func(paths map[string]string) []entity.FileStrategy {
					return []entity.FileStrategy{MockFileProc{}}
				},
				fileExecutorFn: func(producers []entity.FileProducer, strategies []entity.FileStrategy) entity.Executor {
					a.NotEmpty(producers)
					a.NotEmpty(strategies)
					return MockExecutor{}
				},
				removeAllFn: func(old string) error {
					a.NotEmpty(old)
					return nil
				},
				logger: mockLogger,
			}

			dir, err := str.Apply(tmpDir)
			a.NoError(err)
			a.Equal(tmpDir, dir)
		})
	})
	t.Run("real", func(t *testing.T) {
		WithTempDir(t, func(tmpDir string) {
			var (
				pathTempA = filepath.Join(tmpDir, pathA)
				pathTempB = filepath.Join(tmpDir, pathB)
				pathTempC = filepath.Join(tmpDir, pathC)

				expectedPathB = filepath.Join(tmpDir, expectedSubPathB)
				expectedPathC = filepath.Join(tmpDir, expectedSubPathC)

				a          = assert.New(t)
				mockLogger = MockLogger{
					infof: func(format string, args ...any) {
						a.NotEmpty(format)
					},
				}
			)

			CreateFile(t, pathTempA, dataA)
			CreateFile(t, pathTempB, dataB)
			CreateFile(t, pathTempC, dataC)

			str := NewFileSystemStrategy(templateData, nil, nil, mockLogger)

			dir, err := str.Apply(tmpDir)
			a.NoError(err)
			a.Equal(tmpDir, dir)
			a.FileExists(pathTempA)
			a.FileExists(expectedPathB)
			a.FileExists(expectedPathC)

			AssertFileDataEqual(t, pathTempA, []byte(varTemplateVariableValue))
			AssertFileDataEqual(t, expectedPathB, dataB)
			AssertFileDataEqual(t, expectedPathC, dataC)
		})
	})
}

func Test_DryRunFileSystemStrategy(t *testing.T) {
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
	res, err := NewDryRunFileSystemStrategy(mockLogger).Apply(dir)
	assert.NoError(t, err)
	assert.Equal(t, dir, res)
}
