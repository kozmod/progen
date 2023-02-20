package factory

import (
	"fmt"
	"sort"

	"github.com/kozmod/progen/internal/config"
	"github.com/kozmod/progen/internal/entity"
	"github.com/kozmod/progen/internal/exec"
)

func NewExecutorChain(
	conf config.Config,
	templateData map[string]any,
	templateOptions []string,
	logger entity.Logger,
	preprocess,
	dryRun bool,
) (entity.Executor, error) {

	type (
		ExecutorBuilder struct {
			line   int32
			procFn func() (entity.Executor, error)
		}
	)

	builders := make([]ExecutorBuilder, 0, len(conf.Dirs)+len(conf.Files)+len(conf.Cmd)+len(conf.FS))
	for _, dirs := range conf.Dirs {
		d := dirs
		builders = append(builders,
			ExecutorBuilder{
				line: d.Line,
				procFn: func() (entity.Executor, error) {
					return NewMkdirExecutor(d.Val, logger, dryRun)
				},
			})
	}

	var loaders []entity.Preprocessor
	for _, files := range conf.Files {
		f := files
		builders = append(builders,
			ExecutorBuilder{
				line: f.Line,
				procFn: func() (entity.Executor, error) {
					executor, l, err := NewFileExecutor(
						f.Val,
						conf.Settings.HTTP,
						templateData,
						templateOptions,
						logger,
						preprocess,
						dryRun)
					loaders = append(loaders, l...)
					return executor, err
				},
			})
	}

	for _, commands := range conf.Cmd {
		cmd := commands
		builders = append(builders,
			ExecutorBuilder{
				line: cmd.Line,
				procFn: func() (entity.Executor, error) {
					return NewRunCommandExecutor(cmd.Val, logger, dryRun)
				},
			})
	}
	for _, path := range conf.FS {
		fs := path
		builders = append(builders,
			ExecutorBuilder{
				line: fs.Line,
				procFn: func() (entity.Executor, error) {
					return NewFSExecutor(
						fs.Val,
						templateData,
						templateOptions,
						logger,
						dryRun)
				},
			})
	}

	sort.Slice(builders, func(i, j int) bool {
		return builders[i].line < builders[j].line
	})

	executors := make([]entity.Executor, 0, len(builders))
	for i, builder := range builders {
		e, err := builder.procFn()
		if err != nil {
			return nil, fmt.Errorf("configure executor [%d]: %w", i, err)
		}
		if e == nil {
			continue
		}
		executors = append(executors, e)
	}

	return exec.NewPreprocessingChain(loaders, executors), nil
}
