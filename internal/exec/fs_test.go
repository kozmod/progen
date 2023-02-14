package exec

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kozmod/progen/internal/entity"
)

func Test_FileSystemProc(t *testing.T) {
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

		proc := FileSystemProc{
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
			processorsFn: func(paths map[string]string) []entity.FileProc {
				return []entity.FileProc{MockFileProc{}}
			},
			fileExecutorFn: func(producers []entity.FileProducer, processors []entity.FileProc) entity.Executor {
				a.NotEmpty(producers)
				a.NotEmpty(processors)
				return MockExecutor{}
			},
			removeAllFn: func(old string) error {
				a.NotEmpty(old)
				return nil
			},
			logger: mockLogger,
		}

		dir, err := proc.Process(tmpDir)
		a.NoError(err)
		a.Equal(tmpDir, dir)
	})
}
