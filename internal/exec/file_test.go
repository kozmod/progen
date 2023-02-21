package exec

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/kozmod/progen/internal/entity"
)

func Test_PreloadProducer(t *testing.T) {
	t.Parallel()

	var (
		path                = filepath.Join("dir_%d", "name_%d")
		generateDataFilesFn = func(n int) []entity.DataFile {
			files := make([]entity.DataFile, 0, n)
			for i := 0; i < n; i++ {
				files = append(files, entity.DataFile{
					FileInfo: entity.NewFileInfo(fmt.Sprintf(path, i, i)),
				})
			}
			return files
		}

		generateMockProducesFn = func(files []entity.DataFile) []entity.FileProducer {
			producers := make([]entity.FileProducer, 0, len(files))
			for _, file := range files {
				producers = append(producers, &MockProducer{
					file: file,
				})
			}
			return producers
		}
	)

	t.Run("success_and_save_order_v1", func(t *testing.T) {
		var (
			files     = generateDataFilesFn(10)
			producers = generateMockProducesFn(files)

			mockLogger = MockLogger{
				infof: func(format string, args ...any) {
					assert.Equal(t, "file process: %s", format)
					assert.NotEmpty(t, args)
				},
			}
		)

		preloader := NewPreloadProducer(producers, mockLogger)
		err := preloader.Process()
		assert.NoError(t, err)

		for i, exp := range files {
			file, err := preloader.Get()
			assert.NoError(t, err)
			assert.Equalf(t, exp, file, "file_%d", i)
		}
	})

	t.Run("success_and_save_order_v2", func(t *testing.T) {
		files := generateDataFilesFn(10)
		rand.New(rand.NewSource(time.Now().UnixNano())).
			Shuffle(len(files), func(i, j int) {
				files[i], files[j] = files[j], files[i]
			})

		var (
			producers = generateMockProducesFn(files)

			mockLogger = MockLogger{
				infof: func(format string, args ...any) {
					assert.Equal(t, "file process: %s", format)
					assert.NotEmpty(t, args)
				},
			}
		)

		preloader := NewPreloadProducer(producers, mockLogger)
		err := preloader.Process()
		assert.NoError(t, err)

		for i, exp := range files {
			file, err := preloader.Get()
			assert.NoError(t, err)
			assert.Equalf(t, exp, file, "file_%d", i)
		}
	})

	t.Run("error", func(t *testing.T) {
		var (
			expErr   = fmt.Errorf("some_producer_err")
			errIndex = 5

			files     = generateDataFilesFn(10)
			producers = generateMockProducesFn(files)

			mockLogger = MockLogger{
				infof: func(format string, args ...any) {
					assert.Equal(t, "file process: %s", format)
					assert.NotEmpty(t, args)
				},
			}
		)

		errProducer, ok := producers[errIndex].(*MockProducer)
		assert.True(t, ok)
		errProducer.err = expErr

		preloader := NewPreloadProducer(producers, mockLogger)
		err := preloader.Process()
		assert.Error(t, err)
		assert.ErrorIs(t, err, expErr)
	})
}

func Test_TemplateFileStrategy(t *testing.T) {
	t.Parallel()

	const (
		name = "some_file"
		dir  = "/some/path"

		noValue = "<no value>"
	)

	var (
		path          = filepath.Join(dir, name)
		newDataFileFn = func(data string) entity.DataFile {
			return entity.DataFile{
				Data:     []byte(data),
				FileInfo: entity.NewFileInfo(path),
			}
		}
	)

	t.Run("success_exec_template_value", func(t *testing.T) {
		var (
			templateValue = "VAL"
			file          = newDataFileFn(`{{.some.Value}}`)
		)
		str := TemplateFileStrategy{
			templateProcFn: func() entity.TemplateProc {
				return entity.NewTemplateProc(
					map[string]any{"some": map[string]any{"Value": templateValue}},
					nil,
					nil,
				)
			},
		}
		res, err := str.Apply(file)
		assert.NoError(t, err)
		assert.Equal(t, templateValue, string(res.Data))
		assert.Equal(t, file.Name(), res.Name())
		assert.Equal(t, file.Dir(), res.Dir())
	})

	t.Run("success_exec_template_functions", func(t *testing.T) {
		var (
			templateValue = "VAL"
			file          = newDataFileFn(`{{ fn }}`)
		)
		str := TemplateFileStrategy{
			templateProcFn: func() entity.TemplateProc {
				return entity.NewTemplateProc(
					nil,
					map[string]any{
						"fn": func() any { return templateValue },
					},
					nil,
				)
			},
		}
		res, err := str.Apply(file)
		assert.NoError(t, err)
		assert.Equal(t, templateValue, string(res.Data))
		assert.Equal(t, file.Name(), res.Name())
		assert.Equal(t, file.Dir(), res.Dir())
	})
	t.Run("missingkey", func(t *testing.T) {
		t.Run("error", func(t *testing.T) {
			var (
				file = newDataFileFn(`{{ .vars.Some }}`)
			)
			str := TemplateFileStrategy{
				templateProcFn: func() entity.TemplateProc {
					return entity.NewTemplateProc(
						nil,
						nil,
						[]string{fmt.Sprintf("%v=%v", entity.TemplateOptionsMissingKey, entity.MissingKeyError)},
					)
				},
			}
			_, err := str.Apply(file)
			assert.Error(t, err)
		})
		t.Run("default", func(t *testing.T) {
			var (
				file = newDataFileFn(`{{ .vars.Some }}`)
			)
			str := TemplateFileStrategy{
				templateProcFn: func() entity.TemplateProc {
					return entity.NewTemplateProc(
						nil,
						nil,
						[]string{fmt.Sprintf("%v=%v", entity.TemplateOptionsMissingKey, entity.MissingKeyDefault)},
					)
				},
			}
			res, err := str.Apply(file)
			assert.NoError(t, err)
			assert.Equal(t, noValue, string(res.Data))
			assert.Equal(t, file.Name(), res.Name())
			assert.Equal(t, file.Dir(), res.Dir())
		})
	})
}

func Test_ReplacePathFileStrategy(t *testing.T) {
	t.Parallel()

	const (
		oldPathA = "old_1/old_A.go"
		oldPathB = "old_1/old_B.go"
		newPathA = "old_1/new_A.go"
		someData = "DATA_!!!"
	)

	var (
		paths = map[string]string{
			oldPathA: newPathA,
		}
	)

	t.Run("success_replace_old_path_to_new_one", func(t *testing.T) {
		str := ReplacePathFileStrategy{paths: paths}
		res, err := str.Apply(entity.DataFile{Data: []byte(someData), FileInfo: entity.NewFileInfo(oldPathA)})
		assert.NoError(t, err)
		assert.Equal(t, newPathA, res.Path())
		assert.Equal(t, []byte(someData), res.Data)
	})
	t.Run("success_nit_replace_old_path_when_new_one_not_present", func(t *testing.T) {
		str := ReplacePathFileStrategy{paths: paths}
		res, err := str.Apply(entity.DataFile{Data: []byte(someData), FileInfo: entity.NewFileInfo(oldPathB)})
		assert.NoError(t, err)
		assert.Equal(t, oldPathB, res.Path())
		assert.Equal(t, []byte(someData), res.Data)
	})
}

func Test_SaveFileStrategy(t *testing.T) {
	SkipSLowTest(t)

	const (
		someDir  = "some_dir"
		someFile = "some_file.txt"
	)

	WithTempDir(t, func(tmpDir string) {
		var (
			a       = assert.New(t)
			expPath = filepath.Join(tmpDir, someDir, someFile)
			in      = entity.DataFile{
				FileInfo: entity.NewFileInfo(expPath),
				Data:     []byte(`some_file_data`),
			}

			mockLogger = MockLogger{
				infof: func(format string, args ...any) {
					a.NotEmpty(format)
					a.ElementsMatch([]string{expPath}, args)
				},
			}
		)

		res, err := NewSaveFileStrategy(mockLogger).Apply(in)
		a.NoError(err)
		a.Equal(in.Dir(), res.Dir())
		a.Equal(in.Name(), res.Name())
		a.Equal(in.Path(), res.Path())
		a.Equal(in.Data, res.Data)
		a.FileExists(res.Path())

		resData, err := os.ReadFile(res.Path())
		a.NoError(err)
		a.Equal(in.Data, resData)
		a.Equal(res.Data, resData)
	})
}
