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
	logger entity.Logger,
	preprocess,
	dryRun bool,
	templateOptions []string,
) (entity.Executor, error) {

	type (
		ProcGenerator struct {
			line   int32
			procFn func() (entity.Executor, error)
		}
	)

	var generators []ProcGenerator
	for _, dirs := range conf.Dirs {
		d := dirs
		generators = append(generators,
			ProcGenerator{
				line: d.Line,
				procFn: func() (entity.Executor, error) {
					return NewMkdirExecutor(d.Val, logger, dryRun)
				},
			})
	}

	var loaders []entity.Preprocessor
	for _, files := range conf.Files {
		f := files
		generators = append(generators,
			ProcGenerator{
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
		generators = append(generators,
			ProcGenerator{
				line: cmd.Line,
				procFn: func() (entity.Executor, error) {
					return NewRunCommandExecutor(cmd.Val, logger, dryRun)
				},
			})
	}
	for _, path := range conf.FS {
		fs := path
		generators = append(generators,
			ProcGenerator{
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
	for _, exec := range conf.Scripts {
		scripts := exec
		generators = append(generators,
			ProcGenerator{
				line: scripts.Line,
				procFn: func() (entity.Executor, error) {
					return NewRunScriptsExecutor(scripts.Val, logger, dryRun)
				},
			})
	}

	sort.Slice(generators, func(i, j int) bool {
		return generators[i].line < generators[j].line
	})

	processors := make([]entity.Executor, 0, len(generators))
	for i, generator := range generators {
		p, err := generator.procFn()
		if err != nil {
			return nil, fmt.Errorf("configure proc [%d]: %w", i, err)
		}
		if p == nil {
			continue
		}
		processors = append(processors, p)
	}

	return exec.NewPreprocessingChain(loaders, processors), nil
}
