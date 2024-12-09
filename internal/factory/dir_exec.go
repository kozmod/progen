package factory

import (
	"slices"

	"github.com/kozmod/progen/internal/entity"
	"github.com/kozmod/progen/internal/exec"
)

func NewMkdirExecutor(dirs []string, logger entity.Logger, dryRun bool) (entity.Executor, error) {
	if len(dirs) == 0 {
		logger.Infof("mkdir executor: `dir` section is empty")
		return nil, nil
	}

	dirSet := slices.Compact(dirs)

	if dryRun {
		return exec.NewDirExecutor(dirSet, []entity.DirStrategy{exec.NewDryRunMkdirAllStrategy(logger)}), nil
	}

	return exec.NewDirExecutor(dirSet, []entity.DirStrategy{exec.NewMkdirAllStrategy(logger)}), nil
}
