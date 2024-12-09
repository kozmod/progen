package factory

import (
	"slices"

	"github.com/kozmod/progen/internal/entity"
	"github.com/kozmod/progen/internal/exec"
)

func NewRmExecutor(paths []string, logger entity.Logger, dryRun bool) (entity.Executor, error) {
	if len(paths) == 0 {
		logger.Infof("rm executor: `rm` section is empty")
		return nil, nil
	}

	pathsSet := slices.Compact(paths)

	if dryRun {
		return exec.NewRmAllExecutor(pathsSet, []entity.RmStrategy{exec.NewDryRmAllStrategy(logger)}), nil
	}

	return exec.NewRmAllExecutor(pathsSet, []entity.RmStrategy{exec.NewRmAllStrategy(logger)}), nil
}
