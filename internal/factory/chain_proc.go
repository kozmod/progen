package factory

import (
	"fmt"
	"sort"

	"github.com/kozmod/progen/internal/config"
	"github.com/kozmod/progen/internal/entity"
	"github.com/kozmod/progen/internal/proc"
)

func NewProcChain(conf config.Config, order map[string]int, logger entity.Logger) (*proc.Chain, error) {

	type (
		IndexedProc struct {
			index int
			proc  proc.Proc
		}
	)

	confFns := map[string]func(config config.Config) (proc.Proc, error){
		config.TagDirs: func(config config.Config) (proc.Proc, error) {
			return NewMkdirProc(config, logger)
		},
		config.TagFiles: func(config config.Config) (proc.Proc, error) {
			return NewFileProc(conf, logger)
		},
		config.TagCmd: func(config config.Config) (proc.Proc, error) {
			return NewRunCommandProc(config, logger)
		},
	}

	//goland:noinspection SpellCheckingInspection
	indexedProcs := make([]IndexedProc, 0, len(confFns))

	for key, fn := range confFns {
		index, ok := order[key]
		if !ok {
			continue
		}

		configuredProc, err := fn(conf)
		if err != nil {
			return nil, fmt.Errorf("configure proc for [%s]: %w", key, err)
		}

		if configuredProc == nil {
			continue
		}

		indexedProcs = append(indexedProcs, IndexedProc{
			index: index,
			proc:  configuredProc,
		})
	}

	sort.Slice(indexedProcs, func(i, j int) bool {
		return indexedProcs[i].index < indexedProcs[j].index
	})

	//goland:noinspection SpellCheckingInspection
	procs := make([]proc.Proc, 0, len(indexedProcs))
	for _, indexedProc := range indexedProcs {
		procs = append(procs, indexedProc.proc)
	}

	return proc.NewProcChane(procs), nil
}
