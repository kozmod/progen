package factory

import (
	"sort"

	"golang.org/x/xerrors"

	"github.com/kozmod/progen/internal/config"
	"github.com/kozmod/progen/internal/entity"
	"github.com/kozmod/progen/internal/exec"
)

func NewExecutorChain(
	conf config.Config,
	actionFilter actionFilter,
	templateData map[string]any,
	templateOptions []string,
	logger entity.Logger,
	preprocess,
	dryRun bool,
) (entity.Executor, error) {

	type (
		ExecutorBuilder struct {
			action string
			line   int32
			procFn func() (entity.Executor, error)
		}
	)

	builders := make([]ExecutorBuilder, 0, len(conf.Dirs)+len(conf.Files)+len(conf.Cmd)+len(conf.FS))
	for _, dirs := range conf.Dirs {
		var (
			d      = dirs
			action = d.Tag
		)

		if !actionFilter.MatchString(action) {
			continue
		}
		builders = append(builders,
			ExecutorBuilder{
				action: action,
				line:   d.Line,
				procFn: func() (entity.Executor, error) {
					return NewMkdirExecutor(d.Val, logger, dryRun)
				},
			})
	}

	var preprocessors []entity.Preprocessor
	for _, files := range conf.Files {
		var (
			f      = files
			action = f.Tag
		)
		if !actionFilter.MatchString(action) {
			continue
		}
		builders = append(builders,
			ExecutorBuilder{
				action: action,
				line:   f.Line,
				procFn: func() (entity.Executor, error) {
					executor, l, err := NewFileExecutor(
						f.Val,
						conf.Settings.HTTP,
						templateData,
						templateOptions,
						logger,
						preprocess,
						dryRun)
					preprocessors = append(preprocessors, l...)
					return executor, err
				},
			})
	}

	for _, commands := range conf.Cmd {
		var (
			cmd    = commands
			action = cmd.Tag
		)
		if !actionFilter.MatchString(action) {
			continue
		}
		builders = append(builders,
			ExecutorBuilder{
				action: action,
				line:   cmd.Line,
				procFn: func() (entity.Executor, error) {
					return NewRunCommandExecutor(cmd.Val, logger, dryRun)
				},
			})
	}
	for _, path := range conf.FS {
		var (
			fs     = path
			action = fs.Tag
		)
		if actionFilter.MatchString(action) {
			continue
		}
		builders = append(builders,
			ExecutorBuilder{
				action: action,
				line:   fs.Line,
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
	for _, builder := range builders {
		e, err := builder.procFn()
		if err != nil {
			return nil, xerrors.Errorf("configure executor [%s]: %w", builder.action, err)
		}
		if e == nil {
			continue
		}
		executors = append(executors, e)
	}

	return exec.NewPreprocessingChain(preprocessors, executors), nil
}
