package factory

import (
	"slices"

	"github.com/kozmod/progen/internal/entity"
	"github.com/kozmod/progen/internal/exec"
)

type FsModifyExecFactory struct {
	templateData    map[string]any
	templateOptions []string
}

func NewFsModifyExecFactory(
	templateData map[string]any,
	templateOptions []string,
) *FsModifyExecFactory {
	return &FsModifyExecFactory{
		templateData:    templateData,
		templateOptions: templateOptions,
	}
}

func (f FsModifyExecFactory) Create(
	dirs []string,
	logger entity.Logger,
	dryRun bool,
) (entity.Executor, error) {
	if len(dirs) == 0 {
		logger.Infof("fs executor: `dir` section is empty")
		return nil, nil
	}

	dirSet := slices.Compact(dirs)

	if dryRun {
		return exec.NewDirExecutor(dirSet, []entity.DirStrategy{exec.NewDryRunFileSystemModifyStrategy(logger)}), nil
	}

	return exec.NewDirExecutor(dirSet, []entity.DirStrategy{
		exec.NewFileSystemModifyStrategy(
			f.templateData,
			entity.TemplateFnsMap,
			f.templateOptions,
			logger),
	}), nil
}
