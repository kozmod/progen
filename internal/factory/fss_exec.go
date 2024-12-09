package factory

import (
	"github.com/kozmod/progen/internal/entity"
	"github.com/kozmod/progen/internal/exec"
)

type FsSaveExecFactory struct {
	templateData    map[string]any
	templateOptions []string
}

func NewFsSaveExecFactory(
	templateData map[string]any,
	templateOptions []string,
) *FsSaveExecFactory {
	return &FsSaveExecFactory{
		templateData:    templateData,
		templateOptions: templateOptions,
	}
}

func (f FsSaveExecFactory) Create(
	fsList []entity.TargetFs,
	logger entity.Logger,
	dryRun bool,
) (entity.Executor, error) {
	if len(fsList) == 0 {
		logger.Infof("fs executor: `fs save` section is empty")
		return nil, nil
	}

	fsStrategyBydDir := make(map[string][]entity.DirStrategy, len(fsList))
	for _, targetFs := range fsList {
		fsStrategyBydDir[targetFs.TargetDir] = append(
			fsStrategyBydDir[targetFs.TargetDir],
			exec.NewFileSystemSaveStrategy(
				targetFs.Fs,
				f.templateData,
				entity.TemplateFnsMap,
				f.templateOptions,
				logger),
		)
	}

	executors := make([]entity.Executor, 0, len(fsStrategyBydDir))
	for dir, strategy := range fsStrategyBydDir {
		dirs := []string{dir}
		if dryRun {
			executors = append(executors,
				exec.NewDirExecutor(dirs, []entity.DirStrategy{exec.NewDryRunFileSystemSaveStrategy(logger)}),
			)
			continue
		}
		executors = append(executors,
			exec.NewDirExecutor(
				dirs,
				strategy,
			),
		)

	}

	return exec.NewChain(executors), nil
}
