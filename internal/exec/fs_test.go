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

	WithTempDir(t, func(tmpDir string) {
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

			pathTempA = filepath.Join(tmpDir, pathA)
			pathTempB = filepath.Join(tmpDir, pathB)
			pathTempC = filepath.Join(tmpDir, pathC)

			expectedPathB = filepath.Join(tmpDir, fmt.Sprintf("%s/some_file.txt", varTemplateVariableValue))
			expectedPathC = filepath.Join(tmpDir, fmt.Sprintf("c/%s.txt", varTemplateVariableValue))

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

		proc := FileSystemStrategy{
			templateProcFn: func() entity.TemplateProc {
				return entity.NewTemplateProc(map[string]any{
					"var": varTemplateVariableValue,
				},
					nil,
					nil)
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

		dir, err := proc.Apply(tmpDir)
		a.NoError(err)
		a.Equal(tmpDir, dir)
	})
}
