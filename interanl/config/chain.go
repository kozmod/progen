package config

import (
	"fmt"
	"sort"

	"github.com/kozmod/progen/interanl/proc"
)

func MustConfigureProcChain(conf Config, order map[string]int) (*proc.Chain, error) {

	type (
		IndexedProc struct {
			index int
			proc  proc.Proc
		}
	)

	confFns := map[string]func(config Config) (proc.Proc, error){
		tagDirs:  MustConfigureMkdirProc,
		tagFiles: MustConfigureWriteFileProc,
		tagCmd:   MustConfigureRunCommandProc,
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
			return nil, fmt.Errorf("configur proc for [%s]: %w", key, err)
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
