package factory

import (
	"github.com/kozmod/progen/internal/config"
	"github.com/kozmod/progen/internal/entity"
	"github.com/kozmod/progen/internal/proc"
)

func NewMkdirProc(conf config.Config, logger entity.Logger) (proc.Proc, error) {
	if len(conf.Dirs) == 0 {
		return nil, nil
	}

	dirSet := uniqueVal(conf.Dirs)
	return proc.NewMkdirAllProc(dirSet, logger), nil
}

func uniqueVal[T comparable](in []T) []T {
	set := make(map[T]struct{}, len(in))
	out := make([]T, 0, len(in))
	for _, val := range in {
		_, ok := set[val]
		if ok {
			continue
		}
		out = append(out, val)
		set[val] = struct{}{}
	}
	return out
}
