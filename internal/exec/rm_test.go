package exec

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kozmod/progen/internal/entity"
)

func Test_RmAllStrategy(t *testing.T) {
	SkipSLowTest(t)

	const (
		someDir  = "some_dir"
		someFile = "file_name.txt"
	)

	t.Run("rm_dir", func(t *testing.T) {
		WithTempDir(t, func(tmpDir string) {
			var (
				a          = assert.New(t)
				path       = filepath.Join(tmpDir, someDir)
				mockLogger = MockLogger{
					infof: func(format string, args ...any) {
						assert.NotEmpty(t, format)
						assert.ElementsMatch(t, []string{path}, args)
					},
				}
			)

			err := os.MkdirAll(path, os.ModePerm)
			a.NoError(err)
			a.DirExists(path)

			err = NewRmAllStrategy(mockLogger).Apply(path)
			a.NoError(err)
			a.NoDirExists(path)
		})
	})

	t.Run("rm_file", func(t *testing.T) {
		WithTempDir(t, func(tmpDir string) {
			var (
				a          = assert.New(t)
				dir        = filepath.Join(tmpDir, someDir, someFile)
				filePath   = filepath.Join(dir, someFile)
				mockLogger = MockLogger{
					infof: func(format string, args ...any) {
						assert.NotEmpty(t, format)
						assert.ElementsMatch(t, []string{filePath}, args)
					},
				}
			)

			err := os.MkdirAll(dir, os.ModePerm)
			a.NoError(err)
			a.DirExists(dir)

			file, err := os.Create(filePath)
			a.NoError(err)
			a.FileExists(filePath)
			a.Equal(filePath, file.Name())

			err = NewRmAllStrategy(mockLogger).Apply(filePath)
			a.NoError(err)
			a.NoFileExists(filePath)
			a.DirExists(dir)
		})
	})
	t.Run("rm_all_files", func(t *testing.T) {
		WithTempDir(t, func(tmpDir string) {
			const (
				filesQuantity = 5
			)
			var (
				a          = assert.New(t)
				dir        = filepath.Join(tmpDir, someDir)
				rmPath     = filepath.Join(dir, entity.Astrix)
				filesPath  = make([]string, 0, filesQuantity)
				mockLogger = MockLogger{
					infof: func(format string, args ...any) {
						assert.NotEmpty(t, format)
					},
				}
			)

			err := os.MkdirAll(dir, os.ModePerm)
			a.NoError(err)
			a.DirExists(dir)

			for i := 0; i < filesQuantity; i++ {
				var (
					fileExt  = filepath.Ext(someFile)
					fileName = fmt.Sprintf("%s_%d%s", strings.TrimSuffix(someFile, fileExt), i, fileExt)
					filePath = filepath.Join(tmpDir, someDir, fileName)
				)
				filesPath = append(filesPath, filePath)
				_, err = os.Create(filePath)
				a.NoError(err)
				a.FileExists(filePath)
			}

			err = NewRmAllStrategy(mockLogger).Apply(rmPath)
			a.NoError(err)
			a.DirExists(dir)
			for _, path := range filesPath {
				a.NoFileExists(path)
			}
		})
	})

}
