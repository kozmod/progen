package factory

import (
	"github.com/kozmod/progen/internal/entity"
	"github.com/kozmod/progen/internal/exec"
)

func NewFSExecutor(
	dirs []string,
	templateData map[string]any,
	templateOptions []string,
	logger entity.Logger,
	dryRun bool) (entity.Executor, error) {
	if len(dirs) == 0 {
		logger.Infof("fs executor: `dir` section is empty")
		return nil, nil
	}

	dirSet := entity.Unique(dirs)

	if dryRun {
		return exec.NewDirExecutor(dirSet, []entity.DirStrategy{exec.NewDryRunFileSystemStrategy(logger)}), nil
	}

	return exec.NewDirExecutor(dirSet, []entity.DirStrategy{
		exec.NewFileSystemStrategy(
			templateData,
			entity.TemplateFnsMap,
			templateOptions,
			logger),
	}), nil
}
