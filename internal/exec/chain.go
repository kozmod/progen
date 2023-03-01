package exec

import (
	"golang.org/x/xerrors"

	"github.com/kozmod/progen/internal/entity"
)

type PreprocessingChain struct {
	loaders    []entity.Preprocessor
	processors []entity.Executor
}

func NewPreprocessingChain(loaders []entity.Preprocessor, processors []entity.Executor) *PreprocessingChain {
	return &PreprocessingChain{
		loaders:    loaders,
		processors: processors,
	}
}

func (c *PreprocessingChain) Exec() error {
	for i, loader := range c.loaders {
		err := loader.Process()
		if err != nil {
			return xerrors.Errorf("preload [%d]: %w", i, err)
		}
	}

	for i, processor := range c.processors {
		err := processor.Exec()
		if err != nil {
			return xerrors.Errorf("execute proc [%d]: %w", i, err)
		}
	}
	return nil
}
