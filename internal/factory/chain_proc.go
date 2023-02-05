package factory

import (
	"fmt"
	"sort"

	"github.com/kozmod/progen/internal/config"
	"github.com/kozmod/progen/internal/entity"
	"github.com/kozmod/progen/internal/proc"
)

func NewExecutorChain(
	conf config.Config,
	templateData map[string]any,
	logger entity.Logger,
	preload,
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
					executor, l, err := NewFileExecutor(f.Val, conf.Settings.HTTP, templateData, logger, preload, dryRun, templateOptions)
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

	if len(loaders) != 0 {
		return proc.NewPreloadChain(loaders, processors), nil
	}

	return proc.NewProcChain(processors), nil
}
