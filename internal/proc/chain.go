package proc

import "fmt"

type Proc interface {
	Exec() error
}

type PreProcessChain struct {
	processors []Proc
	chain      Proc
}

func NewPreProcessChain(chain Proc, processors ...Proc) *PreProcessChain {
	return &PreProcessChain{
		processors: processors,
		chain:      chain,
	}
}

func (c *PreProcessChain) Exec() error {
	for i, processor := range c.processors {
		err := processor.Exec()
		if err != nil {
			return fmt.Errorf("preprocess chain: preprocessor [%d]: %w", i, err)
		}
	}

	err := c.chain.Exec()
	if err != nil {
		return fmt.Errorf("preprocess chain: %w", err)
	}
	return nil
}

type Chain struct {
	processors []Proc
}

func NewProcChain(processors []Proc) *Chain {
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
