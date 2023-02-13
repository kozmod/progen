package exec

import (
	"fmt"
	"os"

	"github.com/kozmod/progen/internal/entity"
)

type DirExecutor struct {
	dirs       []string
	processors []entity.DirProc
}

func NewDirExecutor(dirs []string, processors []entity.DirProc) *DirExecutor {
	return &DirExecutor{
		dirs:       dirs,
		processors: processors,
	}
}

func (p *DirExecutor) Exec() error {
	for _, dir := range p.dirs {
		path := dir
		for _, processor := range p.processors {
			_, err := processor.Process(path)
			if err != nil {
				return fmt.Errorf("execute dir: process dir [%s]: %w", path, err)
			}
		}
	}
	return nil
}

type MkdirAllProc struct {
	fileMode os.FileMode
	logger   entity.Logger
}

func NewMkdirAllProc(logger entity.Logger) *MkdirAllProc {
	return &MkdirAllProc{
		fileMode: os.ModePerm,
		logger:   logger,
	}
}

func (p *MkdirAllProc) Process(dir string) (string, error) {
	err := os.MkdirAll(dir, p.fileMode)
	if err != nil {
		return entity.Empty, fmt.Errorf("create dir [%s]: %w", dir, err)
	}
	p.logger.Infof("dir created: %s", dir)
	return dir, nil
}

type DryRunMkdirAllProc struct {
	fileMode os.FileMode
	logger   entity.Logger
}

func NewDryRunMkdirAllProc(logger entity.Logger) *DryRunMkdirAllProc {
	return &DryRunMkdirAllProc{
		fileMode: os.ModePerm,
		logger:   logger,
	}
}

func (p *DryRunMkdirAllProc) Process(dir string) (string, error) {
	p.logger.Infof("dir created: %s", dir)
	return dir, nil
}
