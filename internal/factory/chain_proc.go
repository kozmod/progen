package factory

import (
	"fmt"
	"sort"

	"github.com/kozmod/progen/internal/config"
	"github.com/kozmod/progen/internal/entity"
	"github.com/kozmod/progen/internal/proc"
)

func NewProcChain(
	conf config.Config,
	templateData map[string]any,
	logger entity.Logger,
	dryRun bool,
) (*proc.Chain, error) {

	type (
		ProcGenerator struct {
			line   int32
			procFn func() (proc.Proc, error)
		}
	)

	var generators []ProcGenerator
	for _, dirs := range conf.Dirs {
		d := dirs
		generators = append(generators,
			ProcGenerator{
				line: d.Line,
				procFn: func() (proc.Proc, error) {
					return NewMkdirProc(d.Val, logger, dryRun)
				},
			})
	}

	for _, files := range conf.Files {
		f := files
		generators = append(generators,
			ProcGenerator{
				line: f.Line,
				procFn: func() (proc.Proc, error) {
					return NewFileProc(f.Val, conf.HTTP, templateData, logger, dryRun)
				},
			})
	}

	for _, commands := range conf.Cmd {
		cmd := commands
		generators = append(generators,
			ProcGenerator{
				line: cmd.Line,
				procFn: func() (proc.Proc, error) {
					return NewRunCommandProc(cmd.Val, logger, dryRun)
				},
			})
	}

	sort.Slice(generators, func(i, j int) bool {
		return generators[i].line < generators[j].line
	})

	//goland:noinspection SpellCheckingInspection
	procs := make([]proc.Proc, 0, len(generators))
	for i, generator := range generators {
		p, err := generator.procFn()
		if err != nil {
			return nil, fmt.Errorf("configure proc [%d]: %w", i, err)
		}
		if p == nil {
			continue
		}
		procs = append(procs, p)
	}

	return proc.NewProcChain(procs), nil
}
