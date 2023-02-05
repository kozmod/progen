package proc

import (
	"fmt"

	"github.com/kozmod/progen/internal/entity"
)

type PreloadChain struct {
	loaders    []entity.Preprocessor
	processors []entity.Executor
}

func NewPreloadChain(loaders []entity.Preprocessor, processors []entity.Executor) *PreloadChain {
	return &PreloadChain{
		loaders:    loaders,
		processors: processors,
	}
}

func (c *PreloadChain) Exec() error {
	for i, loader := range c.loaders {
		err := loader.Process()
		if err != nil {
			return fmt.Errorf("prealod [%d]: %w", i, err)
		}
	}

	for i, processor := range c.processors {
		err := processor.Exec()
		if err != nil {
			return fmt.Errorf("execute proc [%d]: %w", i, err)
		}
	}
	return nil
}

type Chain struct {
	processors []entity.Executor
}

func NewProcChain(processors []entity.Executor) *Chain {
	return &Chain{
		processors: processors,
	}
}

func (c *Chain) Exec() error {
	for i, processor := range c.processors {
		err := processor.Exec()
		if err != nil {
			return fmt.Errorf("execute proc [%d]: %w", i, err)
		}
	}
	return nil
}
