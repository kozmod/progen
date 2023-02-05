package factory

import (
	"github.com/kozmod/progen/internal/entity"
	"github.com/kozmod/progen/internal/proc"
)

func NewMkdirExecutor(dirs []string, logger entity.Logger, dryRun bool) (entity.Executor, error) {
	if len(dirs) == 0 {
		logger.Infof("`dir` section is empty")
		return nil, nil
	}

	dirSet := entity.Unique(dirs)

	if dryRun {
		return proc.NewDryRunMkdirAllProc(dirSet, logger), nil
	}

	return proc.NewMkdirAllProc(dirSet, logger), nil
}
