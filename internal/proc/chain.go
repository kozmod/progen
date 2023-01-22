package proc

import "fmt"

type Proc interface {
	Exec() error
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
