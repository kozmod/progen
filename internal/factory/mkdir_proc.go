package factory

import (
	"os"

	"github.com/kozmod/progen/internal/config"
	"github.com/kozmod/progen/internal/entity"
	"github.com/kozmod/progen/internal/proc"
)

func NewMkdirProc(conf config.Config, logger entity.Logger, dryRun bool) (proc.Proc, error) {
	if len(conf.Dirs) == 0 {
		return nil, nil
	}

	dirSet := uniqueVal[config.Dir, string](conf.Dirs, func(dir config.Dir) string {
		return dir.Path
	})

	dirs := make([]entity.Dir, 0, len(dirSet))
	for _, dir := range dirSet {
		perm := os.ModePerm
		if dir.Perm != nil {
			perm = dir.Perm.FileMode
		}

		dirs = append(dirs, entity.Dir{Path: dir.Path, Perm: perm})
	}

	if dryRun {
		return proc.NewDryRunMkdirAllProc(dirs, logger), nil
	}

	return proc.NewMkdirAllProc(dirs, logger), nil
}

func uniqueVal[T any, R comparable](in []T, keyFn func(T) R) []T {
	set := make(map[R]struct{}, len(in))
	out := make([]T, 0, len(in))
	for _, val := range in {
		key := keyFn(val)
		_, ok := set[key]
		if ok {
			continue
		}
		out = append(out, val)
		set[key] = struct{}{}
	}
	return out
}
