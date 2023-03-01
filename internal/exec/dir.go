package exec

import (
	"os"

	"golang.org/x/xerrors"

	"github.com/kozmod/progen/internal/entity"
)

type DirExecutor struct {
	dirs       []string
	strategies []entity.DirStrategy
}

func NewDirExecutor(dirs []string, strategies []entity.DirStrategy) *DirExecutor {
	return &DirExecutor{
		dirs:       dirs,
		strategies: strategies,
	}
}

func (p *DirExecutor) Exec() error {
	for _, dir := range p.dirs {
		path := dir
		for _, strategy := range p.strategies {
			_, err := strategy.Apply(path)
			if err != nil {
				return xerrors.Errorf("execute dir: process dir [%s]: %w", path, err)
			}
		}
	}
	return nil
}

type MkdirAllStrategy struct {
	fileMode os.FileMode
	logger   entity.Logger
}

func NewMkdirAllStrategy(logger entity.Logger) *MkdirAllStrategy {
	return &MkdirAllStrategy{
		fileMode: os.ModePerm,
		logger:   logger,
	}
}

func (p *MkdirAllStrategy) Apply(dir string) (string, error) {
	err := os.MkdirAll(dir, p.fileMode)
	if err != nil {
		return entity.Empty, xerrors.Errorf("create dir [%s]: %w", dir, err)
	}
	p.logger.Infof("dir created: %s", dir)
	return dir, nil
}

type DryRunMkdirAllStrategy struct {
	fileMode os.FileMode
	logger   entity.Logger
}

func NewDryRunMkdirAllStrategy(logger entity.Logger) *DryRunMkdirAllStrategy {
	return &DryRunMkdirAllStrategy{
		fileMode: os.ModePerm,
		logger:   logger,
	}
}

func (p *DryRunMkdirAllStrategy) Apply(dir string) (string, error) {
	p.logger.Infof("dir created: %s", dir)
	return dir, nil
}
