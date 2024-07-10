package exec

import (
	"golang.org/x/xerrors"
	"sync"

	"github.com/kozmod/progen/internal/entity"
)

type Chain struct {
	executors []entity.Executor
}

func NewChain(executors []entity.Executor) *Chain {
	return &Chain{executors: executors}
}

func (c *Chain) Exec() error {
	for i, executor := range c.executors {
		err := executor.Exec()
		if err != nil {
			return xerrors.Errorf("execute proc [%d]: %w", i, err)
		}
	}
	return nil
}

type PreprocessingChain struct {
	preprocessors *Preprocessors
	chain         *Chain
}

func NewPreprocessingChain(preprocessors *Preprocessors, executors []entity.Executor) *PreprocessingChain {
	return &PreprocessingChain{
		preprocessors: preprocessors,
		chain:         NewChain(executors),
	}
}

func (c *PreprocessingChain) Exec() error {
	for i, preprocessor := range c.preprocessors.Get() {
		err := preprocessor.Process()
		if err != nil {
			return xerrors.Errorf("preprocess [%d]: %w", i, err)
		}
	}
	return c.chain.Exec()
}

type Preprocessors struct {
	mx  sync.RWMutex
	val []entity.Preprocessor
}

func (p *Preprocessors) Add(in ...entity.Preprocessor) {
	if p == nil {
		return
	}

	p.mx.Lock()
	defer p.mx.Unlock()
	p.val = append(p.val, in...)
}

func (p *Preprocessors) Get() []entity.Preprocessor {
	if p == nil {
		return nil
	}

	p.mx.RLock()
	defer p.mx.RUnlock()
	if len(p.val) == 0 {
		return nil
	}
	res := make([]entity.Preprocessor, len(p.val), len(p.val))
	copy(res, p.val)
	return res
}
