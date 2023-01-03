package factory

import (
	"github.com/kozmod/progen/internal/config"
	"github.com/kozmod/progen/internal/proc"
)

func NewMkdirProc(conf config.Config) (proc.Proc, error) {
	if len(conf.Dirs) == 0 {
		return nil, nil
	}
	return proc.NewMkdirAllProc(uniqueVal(conf.Dirs)), nil
}

func uniqueVal(in []string) []string {
	set := make(map[string]struct{}, len(in))
	out := make([]string, 0, len(in))
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
