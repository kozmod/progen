package exec

import (
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/xerrors"

	"github.com/kozmod/progen/internal/entity"
)

type RmAllExecutor struct {
	paths      []string
	strategies []entity.RmStrategy
}

func NewRmAllExecutor(paths []string, strategies []entity.RmStrategy) *RmAllExecutor {
	return &RmAllExecutor{
		paths:      paths,
		strategies: strategies,
	}
}

func (p *RmAllExecutor) Exec() error {
	for _, path := range p.paths {
		for _, strategy := range p.strategies {
			err := strategy.Apply(path)
			if err != nil {
				return xerrors.Errorf("execute rm: process rm [%s]: %w", path, err)
			}
		}
	}
	return nil
}

type RmAllStrategy struct {
	logger entity.Logger
}

func NewRmAllStrategy(logger entity.Logger) *RmAllStrategy {
	return &RmAllStrategy{
		logger: logger,
	}
}

func (p *RmAllStrategy) Apply(path string) error {
	astrixIndex := strings.Index(path, entity.Astrix)
	if astrixIndex == len(path)-1 {
		contents, err := filepath.Glob(path)
		if err != nil {
			return xerrors.Errorf("rm [%s]: get names of the all files: %w", path, err)
		}
		for _, item := range contents {
			err = os.RemoveAll(item)
			if err != nil {
				return xerrors.Errorf("rm content [%s]: %w", item, err)
			}
		}
		p.logger.Infof("rm all: %s", path)
		return nil
	}

	err := os.RemoveAll(path)
	if err != nil {
		return xerrors.Errorf("rm [%s]: %w", path, err)
	}
	p.logger.Infof("rm: %s", path)
	return nil
}

type DryRunRmAllStrategy struct {
	logger entity.Logger
}

func NewDryRmAllStrategy(logger entity.Logger) *DryRunRmAllStrategy {
	return &DryRunRmAllStrategy{
		logger: logger,
	}
}

func (p *DryRunRmAllStrategy) Apply(path string) error {
	p.logger.Infof("rm: %s", path)
	return nil
}
