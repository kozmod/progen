package exec

import (
	"golang.org/x/xerrors"

	"github.com/kozmod/progen/internal/entity"
)

type PreprocessingChain struct {
	preprocessors []entity.Preprocessor
	executors     []entity.Executor
}

func NewPreprocessingChain(preprocessors []entity.Preprocessor, executors []entity.Executor) *PreprocessingChain {
	return &PreprocessingChain{
		preprocessors: preprocessors,
		executors:     executors,
	}
}

func (c *PreprocessingChain) Exec() error {
	for i, preprocessor := range c.preprocessors {
		err := preprocessor.Process()
		if err != nil {
			return xerrors.Errorf("preload [%d]: %w", i, err)
		}
	}

	for i, executor := range c.executors {
		err := executor.Exec()
		if err != nil {
			return xerrors.Errorf("execute proc [%d]: %w", i, err)
		}
	}
	return nil
}
