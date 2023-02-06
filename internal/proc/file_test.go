package proc

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"

	"github.com/kozmod/progen/internal/entity"
)

func Test_PreloadProducer(t *testing.T) {
	t.Parallel()

	var (
		generateDataFilesFn = func(n int) []entity.DataFile {
			files := make([]entity.DataFile, 0, n)
			for i := 0; i < n; i++ {
				files = append(files, entity.DataFile{
					FileInfo: entity.FileInfo{
						Dir:  fmt.Sprintf("dir_%d", i),
						Name: fmt.Sprintf("name_%d", i),
					},
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
