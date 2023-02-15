package factory

import (
	"github.com/kozmod/progen/internal/entity"
	"github.com/kozmod/progen/internal/exec"
)

func NewMkdirExecutor(dirs []string, logger entity.Logger, dryRun bool) (entity.Executor, error) {
	if len(dirs) == 0 {
		logger.Infof("`dir` section is empty")
		return nil, nil
	}

	dirSet := entity.Unique(dirs)

	if dryRun {
		return exec.NewDirExecutor(dirSet, []entity.DirProc{exec.NewDryRunMkdirAllProc(logger)}), nil
	}

	return exec.NewDirExecutor(dirSet, []entity.DirProc{exec.NewMkdirAllProc(logger)}), nil
}
